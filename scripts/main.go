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

// 定義與前端 index.html 格式對齊的結構
type DailyRow struct {
	IPAddress  string `json:"ip_address"`
	Location   string `json:"location"`
	Model      string `json:"model"`
	RecordedAt string `json:"recorded_at"`
	TotalCount int    `json:"total_count"`
	PrintCount int    `json:"print_count"`
	CopyCount  int    `json:"copy_count"`
	ScanTotal  int    `json:"scan_total"`
	FaxCount   int    `json:"fax_count"`
	DailyTotal int    `json:"daily_total"`
	DailyPrint int    `json:"daily_print"`
	DailyCopy  int    `json:"daily_copy"`
	DailyScan  int    `json:"daily_scan"`
	DailyFax   int    `json:"daily_fax"`
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

	// 3. 查詢我們之前建立的視圖
	// 查詢更新後的 daily_stats 視圖
	rows, err := db.Query(`
		SELECT 
			ip_address, location, model, recorded_at::text, 
			total_count, print_count, copy_count, scan_total, fax_count,
			daily_total, daily_print, daily_copy, daily_scan, daily_fax
		FROM daily_stats
		ORDER BY recorded_at DESC
	`)
	if err != nil {
		log.Fatal("查詢失敗:", err)
	}
	defer rows.Close()

	var results []DailyRow
	for rows.Next() {
		var r DailyRow
		err := rows.Scan(
			&r.IPAddress, &r.Location, &r.Model, &r.RecordedAt,
			&r.TotalCount, &r.PrintCount, &r.CopyCount, &r.ScanTotal, &r.FaxCount,
			&r.DailyTotal, &r.DailyPrint, &r.DailyCopy, &r.DailyScan, &r.DailyFax,
		)
		if err != nil {
			log.Fatal("讀取資料失敗:", err)
		}
		results = append(results, r)
	}

	// 4. 將資料包裝成前端預期的格式
	// 我們直接輸出這組陣列，由前端 index.html 中的 transformData 處理
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal("JSON 轉換失敗:", err)
	}

	// 5. 確保輸出目錄存在並寫入檔案
	err = os.MkdirAll("public", 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("public/data.json", jsonData, 0644)
	if err != nil {
		log.Fatal("檔案寫入失敗:", err)
	}

	fmt.Printf("成功！已從 Supabase 擷取 %d 筆紀錄並儲存至 public/data.json\n", len(results))
	fmt.Println("產生時間:", time.Now().Format("2006-01-02 15:04:05"))
}
