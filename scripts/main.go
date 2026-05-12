package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// 每日統計行（來自 daily_stats 視圖）
type DailyRow struct {
	IPAddress  string `json:"ip_address"`
	RecordedAt string `json:"recorded_at"`
	DailyTotal int    `json:"daily_total"`
	TotalCount int    `json:"total_count,omitempty"`
	PrintCount int    `json:"print_count,omitempty"`
	CopyCount  int    `json:"copy_count,omitempty"`
	ScanTotal  int    `json:"scan_total,omitempty"`
	FaxCount   int    `json:"fax_count,omitempty"`
	DailyPrint int    `json:"daily_print,omitempty"`
	DailyCopy  int    `json:"daily_copy,omitempty"`
	DailyScan  int    `json:"daily_scan,omitempty"`
	DailyFax   int    `json:"daily_fax,omitempty"`
}

// 耗材狀態
type Supply struct {
	Type    int `json:"type"`
	Percent int `json:"percent"`
}

// 故障紀錄
type Incident struct {
	ErrorCode    string `json:"error_code"`
	Description  string `json:"description"`
	IncidentDate string `json:"incident_date"`
	DowntimeMin  int    `json:"downtime_minutes"`
}

// 儀表板印表機（單一物件，供前端直接使用）
type DashboardPrinter struct {
	ID              int        `json:"id"`
	IPAddress       string     `json:"ip_address"`
	Location        string     `json:"location"`
	Model           string     `json:"model"`
	Unit            string     `json:"unit"`
	PrinterState    string     `json:"printer_status"`
	TonerPercent    int        `json:"toner_percent"`
	InkPercent      int        `json:"ink_percent"`
	PaperPercent    int        `json:"paper_percent"`
	Supplies        []Supply   `json:"supplies"`
	WarrantyDays    int        `json:"warranty_days"`
	RecentIncidents int        `json:"recent_incidents_30d"`
	Trend           []DailyRow `json:"trend"`
	HistoryLogs     []DailyRow `json:"history_logs"`
	Incidents       []Incident `json:"incidents"`
	TotalCount      int        `json:"total_all_time"`
}

// 整體儀表板回應
type DashboardResponse struct {
	Summary struct {
		TotalPrinters      int `json:"total_printers"`
		TodayTotal         int `json:"today_total"`
		SupplyWarningCount int `json:"supply_warning_count"`
	} `json:"summary"`
	Printers []DashboardPrinter `json:"printers"`
}

func main() {
	// 1. 從環境變數獲取連線字串
	connStr := os.Getenv("SUPABASE_DB_CONNECTION")
	if connStr == "" {
		log.Fatal("錯誤: 未設定 SUPABASE_DB_CONNECTION 環境變數")
	}

	// 2. 連線至資料庫
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. 查詢 daily_stats 視圖（歷史趨勢，包含每日增量 daily_total）
	trendRows, err := db.Query(`
		SELECT
			ip_address, recorded_at::text,
			COALESCE(total_count, 0),
			COALESCE(print_count, 0),
			COALESCE(copy_count, 0),
			COALESCE(scan_total, 0),
			COALESCE(fax_count, 0)
		FROM daily_stats
		ORDER BY ip_address, recorded_at ASC
	`)
	if err != nil {
		log.Fatal("趨勢查詢失敗:", err)
	}

	// 按 IP 分組趨勢資料
	trendMap := make(map[string][]DailyRow)
	for trendRows.Next() {
		var r DailyRow
		err := trendRows.Scan(
			&r.IPAddress, &r.RecordedAt,
			&r.TotalCount,
			&r.PrintCount,
			&r.CopyCount,
			&r.ScanTotal,
			&r.FaxCount,
		)
		if err != nil {
			log.Fatal("趨勢讀取失敗:", err)
		}
		trendMap[r.IPAddress] = append(trendMap[r.IPAddress], r)
	}
	trendRows.Close()
	for ip, rows := range trendMap {
		trendMap[ip] = buildCorrectedTrend(rows)
	}

	// 4. 查詢 dashboard_stats 視圖（當前狀態）
	dashRows, err := db.Query(`
		SELECT id, ip_address, location, model, unit, printer_status,
		       toner_percent, ink_percent, paper_percent, warranty_days, recent_incidents_30d
		FROM dashboard_stats
		ORDER BY unit, ip_address
	`)
	if err != nil {
		log.Fatal("儀表板查詢失敗:", err)
	}

	// 查詢耗材詳細資料
	supplyRows, err := db.Query(`SELECT printer_id, supply_type, remaining_percent FROM supplies`)
	if err != nil {
		log.Fatal("耗材查詢失敗:", err)
	}

	// 查詢近期故障紀錄
	incidentRows, err := db.Query(`
		SELECT printer_id, error_code, description, incident_date::text, downtime_minutes
		FROM incidents
		WHERE incident_date >= CURRENT_DATE - INTERVAL '30 days'
		ORDER BY incident_date DESC
	`)
	if err != nil {
		log.Fatal("故障紀錄查詢失敗:", err)
	}

	// 將 supplyRows 轉換為 map
	supplyMap := make(map[int][]Supply)
	for supplyRows.Next() {
		var printerID int
		var supplyType string
		var remainingPercent int
		err := supplyRows.Scan(&printerID, &supplyType, &remainingPercent)
		if err != nil {
			log.Fatal("耗材讀取失敗:", err)
		}
		supplyMap[printerID] = append(supplyMap[printerID], Supply{
			Type:    supplyTypeToInt(supplyType),
			Percent: remainingPercent,
		})
	}
	supplyRows.Close()

	// 將 incidentRows 轉換為 map
	incidentMap := make(map[int][]Incident)
	for incidentRows.Next() {
		var printerID int
		var inc Incident
		err := incidentRows.Scan(&printerID, &inc.ErrorCode, &inc.Description, &inc.IncidentDate, &inc.DowntimeMin)
		if err != nil {
			log.Fatal("故障紀錄讀取失敗:", err)
		}
		incidentMap[printerID] = append(incidentMap[printerID], inc)
	}
	incidentRows.Close()

	// 組合最終儀表板資料
	response := DashboardResponse{}
	var seenIPs map[string]bool
	response.Printers = make([]DashboardPrinter, 0)
	seenIPs = make(map[string]bool)

	for dashRows.Next() {
		var dp DashboardPrinter
		var warrantyDays sql.NullInt64
		err := dashRows.Scan(
			&dp.ID, &dp.IPAddress, &dp.Location, &dp.Model, &dp.Unit, &dp.PrinterState,
			&dp.TonerPercent, &dp.InkPercent, &dp.PaperPercent,
			&warrantyDays, &dp.RecentIncidents,
		)
		if err != nil {
			log.Fatal("儀表板讀取失敗:", err)
		}
		// 若資料庫中 warranty_days 為 NULL，則轉換為 0
		if warrantyDays.Valid {
			dp.WarrantyDays = int(warrantyDays.Int64)
		} else {
			dp.WarrantyDays = 0
		}

		dp.Supplies = supplyMap[dp.ID]
		dp.Incidents = incidentMap[dp.ID]
		dp.Trend = trendMap[dp.IPAddress]

		// 從趨勢資料累加 total_all_time
		dp.TotalCount = 0
		for _, t := range dp.Trend {
			dp.TotalCount += t.DailyTotal
		}

		// 去重（dashboard 和趨勢可能有重疊的 IP）
		if !seenIPs[dp.IPAddress] {
			seenIPs[dp.IPAddress] = true
			response.Printers = append(response.Printers, dp)
		}
	}
	dashRows.Close()

	// 統計概況
	response.Summary.TotalPrinters = len(response.Printers)
	today := dashboardToday()
	for _, p := range response.Printers {
		for _, s := range p.Supplies {
			if s.Percent < 15 {
				response.Summary.SupplyWarningCount++
			}
		}
		for _, row := range p.Trend {
			if dateOnly(row.RecordedAt) == today && row.DailyTotal > 0 {
				response.Summary.TodayTotal += row.DailyTotal
			}
		}
	}

	// 將資料寫入 JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal("JSON 轉換失敗:", err)
	}

	// 確保輸出目錄存在並寫入檔案
	err = os.MkdirAll("public", 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("public/data.json", jsonData, 0644)
	if err != nil {
		log.Fatal("檔案寫入失敗:", err)
	}

	fmt.Printf("成功！已從 Supabase 擷取 %d 台印表機資料並儲存至 public/data.json\n", len(response.Printers))
	fmt.Println("產生時間:", time.Now().Format("2006-01-02 15:04:05"))
}

func supplyTypeToInt(t string) int {
	switch t {
	case "toner":
		return 0
	case "ink":
		return 1
	case "paper":
		return 2
	default:
		return 3
	}
}

func dashboardToday() string {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return time.Now().Format("2006-01-02")
	}
	return time.Now().In(loc).Format("2006-01-02")
}

func dateOnly(value string) string {
	if len(value) >= len("2006-01-02") {
		return value[:len("2006-01-02")]
	}
	return value
}

func buildCorrectedTrend(rows []DailyRow) []DailyRow {
	var previous *DailyRow
	for i := range rows {
		current := rows[i]
		hasReading := hasCounterReading(current)

		if hasReading && previous != nil {
			rows[i].DailyPrint = counterDelta(current.PrintCount, previous.PrintCount)
			rows[i].DailyCopy = counterDelta(current.CopyCount, previous.CopyCount)
			rows[i].DailyScan = counterDelta(current.ScanTotal, previous.ScanTotal)
			rows[i].DailyFax = counterDelta(current.FaxCount, previous.FaxCount)

			// 總數量 = 列印 + 影印 + 傳真（不含掃描）
			rows[i].DailyTotal = rows[i].DailyPrint + rows[i].DailyCopy + rows[i].DailyFax
		}

		if hasReading {
			previous = &rows[i]
		} else if previous == nil {
			previous = &rows[i]
		}
	}
	return rows
}

func hasCounterReading(row DailyRow) bool {
	return row.TotalCount > 0 ||
		row.PrintCount > 0 ||
		row.CopyCount > 0 ||
		row.ScanTotal > 0 ||
		row.FaxCount > 0
}

func counterDelta(current, previous int) int {
	if current <= previous {
		return 0
	}
	return current - previous
}
