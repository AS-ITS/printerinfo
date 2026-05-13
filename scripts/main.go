package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// API 回應格式
type APIResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	db              *sql.DB
	trendMap        = make(map[string][]DailyRow)
	supplyMap       = make(map[int][]Supply)
	incidentMap     = make(map[int][]Incident)
	dashboardCache  []byte
	cacheMutex      struct{}
	dataLoaded      bool
)

func main() {
	// 1. 從環境變數獲取連線字串
	connStr := os.Getenv("SUPABASE_DB_CONNECTION")
	if connStr == "" {
		log.Fatal("錯誤: 未設定 SUPABASE_DB_CONNECTION 環境變數")
	}

	// 2. 連線至資料庫
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. 初始載入資料
	if err := loadAllData(); err != nil {
		log.Printf("初始載入失敗: %v", err)
	}

	// 4. 定期更新資料
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := loadAllData(); err != nil {
				log.Printf("定期更新失敗: %v", err)
			}
		}
	}()

	// 5. HTTP Server
	mux := http.NewServeMux()

	// API 端點
	mux.HandleFunc("/api/printer", handlerPrinter)
	mux.HandleFunc("/api/printer/", handlerPrinterByID)

	// 靜態檔案
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// loadAllData 從資料庫載入所有資料並更新緩存
func loadAllData() error {
	connStr := os.Getenv("SUPABASE_DB_CONNECTION")
	if connStr == "" {
		return fmt.Errorf("未設定 SUPABASE_DB_CONNECTION")
	}

	// 重新連線
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("趨勢查詢失敗: %w", err)
	}

	// 按 IP 分組趨勢資料
	trendMap = make(map[string][]DailyRow)
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
			return fmt.Errorf("趨勢讀取失敗: %w", err)
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
		return fmt.Errorf("儀表板查詢失敗: %w", err)
	}

	// 查詢耗材詳細資料
	supplyRows, err := db.Query(`SELECT printer_id, supply_type, remaining_percent FROM supplies`)
	if err != nil {
		return fmt.Errorf("耗材查詢失敗: %w", err)
	}

	// 查詢近期故障紀錄
	incidentRows, err := db.Query(`
		SELECT printer_id, error_code, description, incident_date::text, downtime_minutes
		FROM incidents
		WHERE incident_date >= CURRENT_DATE - INTERVAL '30 days'
		ORDER BY incident_date DESC
	`)
	if err != nil {
		return fmt.Errorf("故障紀錄查詢失敗: %w", err)
	}

	// 將 supplyRows 轉換為 map
	supplyMap = make(map[int][]Supply)
	for supplyRows.Next() {
		var printerID int
		var supplyType string
		var remainingPercent int
		err := supplyRows.Scan(&printerID, &supplyType, &remainingPercent)
		if err != nil {
			return fmt.Errorf("耗材讀取失敗: %w", err)
		}
		supplyMap[printerID] = append(supplyMap[printerID], Supply{
			Type:    supplyTypeToInt(supplyType),
			Percent: remainingPercent,
		})
	}
	supplyRows.Close()

	// 將 incidentRows 轉換為 map
	incidentMap = make(map[int][]Incident)
	for incidentRows.Next() {
		var printerID int
		var inc Incident
		err := incidentRows.Scan(&printerID, &inc.ErrorCode, &inc.Description, &inc.IncidentDate, &inc.DowntimeMin)
		if err != nil {
			return fmt.Errorf("故障紀錄讀取失敗: %w", err)
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
			return fmt.Errorf("儀表板讀取失敗: %w", err)
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
		return fmt.Errorf("JSON 轉換失敗: %w", err)
	}

	// 確保輸出目錄存在並寫入檔案
	err = os.MkdirAll("public", 0755)
	if err != nil {
		return fmt.Errorf("目錄創建失敗: %w", err)
	}

	err = os.WriteFile("public/data.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("檔案寫入失敗: %w", err)
	}

	dashboardCache = jsonData
	dataLoaded = true

	fmt.Printf("成功！已從 Supabase 擷取 %d 台印表機資料並儲存至 public/data.json\n", len(response.Printers))
	fmt.Println("產生時間:", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

// handlerPrinter 處理 /api/printer 端點
func handlerPrinter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "Method not allowed"})
		return
	}

	if !dataLoaded {
		if err := loadAllData(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
			return
		}
	}

	// 檢查是否指定單一印表機
	idStr := r.URL.Query().Get("id")
	ip := r.URL.Query().Get("ip")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "id 必須是數字"})
			return
		}
		handlerPrinterByID(w, r)
		return
	}

	if ip != "" {
		// 解析單一印表機
		printer, ok := getPrinterByIP(ip)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "找不到該 IP 的印表機"})
			return
		}
		json.NewEncoder(w).Encode(APIResponse{Success: true, Data: printer})
		return
	}

	// 返回全部
	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: dashboardCache})
}

// handlerPrinterByID 處理 /api/printer/{id} 端點
func handlerPrinterByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "Method not allowed"})
		return
	}

	if !dataLoaded {
		if err := loadAllData(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
			return
		}
	}

	// 從 URL path 提取 ID
	path := strings.TrimPrefix(r.URL.Path, "/api/printer/")
	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "請提供印表機 ID"})
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "ID 必須是數字"})
		return
	}

	// 查找印表機
	for _, p := range dashboardCacheData() {
		if p.ID == id {
			json.NewEncoder(w).Encode(APIResponse{Success: true, Data: p})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Success: false, Error: fmt.Sprintf("找不到 ID 為 %d 的印表機", id)})
}

// dashboardCacheData 解析快取的儀表板資料
func dashboardCacheData() []DashboardPrinter {
	if !dataLoaded || dashboardCache == nil {
		return nil
	}
	var resp DashboardResponse
	if err := json.Unmarshal(dashboardCache, &resp); err != nil {
		return nil
	}
	return resp.Printers
}

// getPrinterByIP 依 IP 地址查找印表機
func getPrinterByIP(ip string) (DashboardPrinter, bool) {
	printers := dashboardCacheData()
	for _, p := range printers {
		if p.IPAddress == ip {
			return p, true
		}
	}
	return DashboardPrinter{}, false
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
