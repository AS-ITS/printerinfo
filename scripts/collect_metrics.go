package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	_ "github.com/lib/pq"
)

const defaultMarkerLifeCountOID = "1.3.6.1.2.1.43.10.2.1.4.1.1"

type inventoryPrinter struct {
	ID        int
	IP        string
	Location  string
	Model     string
	Unit      string
	PrevCount metricCounts
}

type metricCounts struct {
	Total int
	Print int
	Copy  int
	Scan  int
	Fax   int
}

type collectorConfig struct {
	Date      string
	Community string
	Port      uint16
	Timeout   time.Duration
	Retries   int
	PrinterID int
	DryRun    bool
	TotalOID  string
	PrintOID  string
	CopyOID   string
	ScanOID   string
	FaxOID    string
}

func main() {
	cfg := parseConfig()

	connStr := os.Getenv("SUPABASE_DB_CONNECTION")
	if connStr == "" {
		log.Fatal("錯誤: 未設定 SUPABASE_DB_CONNECTION 環境變數")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	printers, err := loadPrinters(db, cfg.PrinterID)
	if err != nil {
		log.Fatal("讀取 printers 失敗:", err)
	}
	if len(printers) == 0 {
		log.Println("沒有可採集的印表機 IP")
		return
	}

	var okCount, failCount int
	for _, printer := range printers {
		counts, err := collectPrinterCounts(printer, cfg)
		if err != nil {
			failCount++
			log.Printf("略過 %s %s (%s): %v", printer.Unit, printer.Location, printer.IP, err)
			continue
		}

		if cfg.DryRun {
			okCount++
			log.Printf("[dry-run] %s %s (%s) total=%d print=%d copy=%d scan=%d fax=%d date=%s",
				printer.Unit, printer.Location, printer.IP,
				counts.Total, counts.Print, counts.Copy, counts.Scan, counts.Fax, cfg.Date)
			continue
		}

		if err := upsertMetric(db, printer.ID, cfg.Date, counts); err != nil {
			failCount++
			log.Printf("寫入失敗 %s %s (%s): %v", printer.Unit, printer.Location, printer.IP, err)
			continue
		}

		okCount++
		log.Printf("已寫入 %s %s (%s) total=%d print=%d copy=%d scan=%d fax=%d date=%s",
			printer.Unit, printer.Location, printer.IP,
			counts.Total, counts.Print, counts.Copy, counts.Scan, counts.Fax, cfg.Date)
	}

	log.Printf("printer_metrics 採集完成: success=%d failed=%d date=%s", okCount, failCount, cfg.Date)
	if okCount == 0 && failCount > 0 {
		os.Exit(1)
	}
}

func parseConfig() collectorConfig {
	defaultDate := time.Now().Format("2006-01-02")
	date := flag.String("date", getenv("PRINTER_METRICS_DATE", defaultDate), "寫入 printer_metrics.recorded_at 的日期，格式 YYYY-MM-DD")
	community := flag.String("community", getenv("SNMP_COMMUNITY", "public"), "SNMP v2c community")
	port := flag.Int("port", getenvInt("SNMP_PORT", 161), "SNMP UDP port")
	timeout := flag.Duration("timeout", time.Duration(getenvInt("SNMP_TIMEOUT_SECONDS", 3))*time.Second, "SNMP timeout，例如 3s")
	retries := flag.Int("retries", getenvInt("SNMP_RETRIES", 1), "SNMP retry 次數")
	printerID := flag.Int("printer-id", getenvInt("PRINTER_ID", 0), "只採集指定 printers.id，0 表示全部")
	dryRun := flag.Bool("dry-run", getenvBool("DRY_RUN", false), "只讀取並列印結果，不寫入資料庫")
	totalOID := flag.String("total-oid", getenv("PRINTER_METRICS_TOTAL_OID", defaultMarkerLifeCountOID), "累積總頁數 OID，預設為 Printer-MIB prtMarkerLifeCount")
	printOID := flag.String("print-oid", getenv("PRINTER_METRICS_PRINT_OID", ""), "列印累積值 OID；未設定時使用 total-oid 的值")
	copyOID := flag.String("copy-oid", getenv("PRINTER_METRICS_COPY_OID", ""), "影印累積值 OID；未設定時沿用前一次資料")
	scanOID := flag.String("scan-oid", getenv("PRINTER_METRICS_SCAN_OID", ""), "掃描累積值 OID；未設定時沿用前一次資料")
	faxOID := flag.String("fax-oid", getenv("PRINTER_METRICS_FAX_OID", ""), "傳真累積值 OID；未設定時沿用前一次資料")
	flag.Parse()

	if _, err := time.Parse("2006-01-02", *date); err != nil {
		log.Fatalf("date 格式錯誤，需為 YYYY-MM-DD: %v", err)
	}
	if *port <= 0 || *port > math.MaxUint16 {
		log.Fatalf("SNMP port 超出範圍: %d", *port)
	}
	if strings.TrimSpace(*totalOID) == "" && strings.TrimSpace(*printOID) == "" {
		log.Fatal("至少需要設定 total-oid 或 print-oid")
	}

	return collectorConfig{
		Date:      *date,
		Community: *community,
		Port:      uint16(*port),
		Timeout:   *timeout,
		Retries:   *retries,
		PrinterID: *printerID,
		DryRun:    *dryRun,
		TotalOID:  normalizeOID(*totalOID),
		PrintOID:  normalizeOID(*printOID),
		CopyOID:   normalizeOID(*copyOID),
		ScanOID:   normalizeOID(*scanOID),
		FaxOID:    normalizeOID(*faxOID),
	}
}

func loadPrinters(db *sql.DB, printerID int) ([]inventoryPrinter, error) {
	query := `
		SELECT id, ip_address, COALESCE(location, ''), COALESCE(model, ''), COALESCE(unit, '')
		FROM printers
		WHERE NULLIF(TRIM(ip_address), '') IS NOT NULL`
	args := []any{}
	if printerID > 0 {
		query += " AND id = $1"
		args = append(args, printerID)
	}
	query += " ORDER BY unit, location, ip_address"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	printers := []inventoryPrinter{}
	for rows.Next() {
		var p inventoryPrinter
		if err := rows.Scan(&p.ID, &p.IP, &p.Location, &p.Model, &p.Unit); err != nil {
			return nil, err
		}
		prev, err := loadPreviousMetric(db, p.ID)
		if err != nil {
			return nil, err
		}
		p.PrevCount = prev
		printers = append(printers, p)
	}
	return printers, rows.Err()
}

func loadPreviousMetric(db *sql.DB, printerID int) (metricCounts, error) {
	var counts metricCounts
	err := db.QueryRow(`
		SELECT
			COALESCE(total_count, 0),
			COALESCE(print_count, 0),
			COALESCE(copy_count, 0),
			COALESCE(scan_total, 0),
			COALESCE(fax_count, 0)
		FROM printer_metrics
		WHERE printer_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1`, printerID).Scan(
		&counts.Total,
		&counts.Print,
		&counts.Copy,
		&counts.Scan,
		&counts.Fax,
	)
	if err == sql.ErrNoRows {
		return metricCounts{}, nil
	}
	return counts, err
}

func collectPrinterCounts(printer inventoryPrinter, cfg collectorConfig) (metricCounts, error) {
	client := &gosnmp.GoSNMP{
		Target:    printer.IP,
		Port:      cfg.Port,
		Community: cfg.Community,
		Version:   gosnmp.Version2c,
		Timeout:   cfg.Timeout,
		Retries:   cfg.Retries,
	}

	if err := client.Connect(); err != nil {
		return metricCounts{}, err
	}
	defer client.Conn.Close()

	counts := printer.PrevCount
	var err error

	if cfg.TotalOID != "" {
		counts.Total, err = getSNMPCounter(client, cfg.TotalOID)
		if err != nil {
			return metricCounts{}, fmt.Errorf("讀取 total OID %s 失敗: %w", cfg.TotalOID, err)
		}
	}

	if cfg.PrintOID != "" {
		counts.Print, err = getSNMPCounter(client, cfg.PrintOID)
		if err != nil {
			return metricCounts{}, fmt.Errorf("讀取 print OID %s 失敗: %w", cfg.PrintOID, err)
		}
	} else if cfg.TotalOID != "" {
		counts.Print = counts.Total
	}

	if cfg.CopyOID != "" {
		counts.Copy, err = getSNMPCounter(client, cfg.CopyOID)
		if err != nil {
			return metricCounts{}, fmt.Errorf("讀取 copy OID %s 失敗: %w", cfg.CopyOID, err)
		}
	}
	if cfg.ScanOID != "" {
		counts.Scan, err = getSNMPCounter(client, cfg.ScanOID)
		if err != nil {
			return metricCounts{}, fmt.Errorf("讀取 scan OID %s 失敗: %w", cfg.ScanOID, err)
		}
	}
	if cfg.FaxOID != "" {
		counts.Fax, err = getSNMPCounter(client, cfg.FaxOID)
		if err != nil {
			return metricCounts{}, fmt.Errorf("讀取 fax OID %s 失敗: %w", cfg.FaxOID, err)
		}
	}
	if counts.Total == 0 {
		counts.Total = counts.Print + counts.Copy + counts.Scan + counts.Fax
	}
	return counts, nil
}

func getSNMPCounter(client *gosnmp.GoSNMP, oid string) (int, error) {
	result, err := client.Get([]string{oid})
	if err != nil {
		return 0, err
	}
	if result == nil || len(result.Variables) == 0 {
		return 0, fmt.Errorf("SNMP 無回傳值")
	}

	pdu := result.Variables[0]
	switch pdu.Type {
	case gosnmp.NoSuchObject, gosnmp.NoSuchInstance, gosnmp.EndOfMibView:
		return 0, fmt.Errorf("OID 不存在或該設備不支援: %s", pdu.Type)
	}

	value := gosnmp.ToBigInt(pdu.Value)
	if value.Sign() < 0 {
		return 0, fmt.Errorf("SNMP counter 為負值: %s", value.String())
	}
	if !value.IsInt64() || value.Int64() > math.MaxInt32 {
		return 0, fmt.Errorf("SNMP counter 超出 PostgreSQL integer 範圍: %s", value.String())
	}
	return int(value.Int64()), nil
}

func upsertMetric(db *sql.DB, printerID int, recordedAt string, counts metricCounts) error {
	_, err := db.Exec(`
		INSERT INTO printer_metrics (
			printer_id,
			total_count,
			print_count,
			copy_count,
			scan_total,
			fax_count,
			recorded_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (printer_id, recorded_at)
		DO UPDATE SET
			total_count = EXCLUDED.total_count,
			print_count = EXCLUDED.print_count,
			copy_count = EXCLUDED.copy_count,
			scan_total = EXCLUDED.scan_total,
			fax_count = EXCLUDED.fax_count`,
		printerID,
		counts.Total,
		counts.Print,
		counts.Copy,
		counts.Scan,
		counts.Fax,
		recordedAt,
	)
	return err
}

func normalizeOID(oid string) string {
	oid = strings.TrimSpace(oid)
	return strings.TrimPrefix(oid, ".")
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
