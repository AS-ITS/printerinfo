package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL 驅動程式
)

// Data結構對應 daily_stats 視圖的欄位
type PrinterStat struct {
	IPAddress  string `json:"ip_address"`
	Location   string `json:"location"`
	Model      string `json:"model"`
	RecordedAt string `json:"recorded_at"`
	TotalCount int    `json:"total_count"`
	PrintCount int    `json:"print_count"`
	CopyCount  int    `json:"copy_count"`
	ScanTotal  int    `json:"scan_total"`
	FaxCount   int    `json:"fax_count"`
	DailyPrint int    `json:"daily_print"`
	DailyCopy  int    `json:"daily_copy"`
	DailyScan  int    `json:"daily_scan"`
	DailyFax   int    `json:"daily_fax"`
}

func main() {
	// 1. 取得環境變數（從 GitHub Secrets 傳入）
	// Supabase 提供標準的連接字串：postgres://postgres:[PASSWORD]@db.[REF].supabase.co:5432/postgres
	connStr := os.Getenv("SUPABASE_DB_CONNECTION")
	if connStr == "" {
		log.Fatal("錯誤：請設定 SUPABASE_DB_CONNECTION 環境變數")
	}

	// 2. 連接到資料庫
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. 執行查詢（讀取 daily_stats 視圖）
	rows, err := db.Query(`
		SELECT 
			ip_address, location, model, recorded_at::text, 
			total_count, print_count, copy_count, scan_total, fax_count,
			daily_print, daily_copy, daily_scan, daily_fax
		FROM daily_stats
		ORDER BY recorded_at DESC
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stats []PrinterStat
	for rows.Next() {
		var s PrinterStat
		err := rows.Scan(
			&s.IPAddress, &s.Location, &s.Model, &s.RecordedAt,
			&s.TotalCount, &s.PrintCount, &s.CopyCount, &s.ScanTotal, &s.FaxCount,
			&s.DailyPrint, &s.DailyCopy, &s.DailyScan, &s.DailyFax,
		)
		if err != nil {
			log.Fatal(err)
		}
		stats = append(stats, s)
	}

	// 4. 將結果轉為 JSON
	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// 5. 確保目錄存在並寫入檔案
	err = os.MkdirAll("public", 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("public/data.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("成功！資料已寫入 public/data.json，總筆數：", len(stats))
}
