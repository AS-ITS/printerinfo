package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// PrinterData 定義與 CSV 欄位對應的結構
type PrinterData struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
	Model    string `json:"model"`
	Total    int    `json:"total"`
	Print    int    `json:"print"`
	Copy     int    `json:"copy"`
	Scan     int    `json:"scan"`
	Date     string `json:"date"`
}

// readCSV 負責讀取單一 CSV 檔案並轉換為 Map，以 IP 為 Key
func readCSV(path string) (map[string]PrinterData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("無法開啟檔案 %s: %v", path, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// 讀取所有內容
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析 CSV 失敗 %s: %v", path, err)
	}

	dataMap := make(map[string]PrinterData)
	for i, r := range records {
		// 跳過首行標題或空行
		if i == 0 || len(r) < 11 {
			continue
		}

		// 嚴謹的數值轉換：若資料異常則設為 0，避免程式崩潰
		total, _ := strconv.Atoi(r[3])
		printVal, _ := strconv.Atoi(r[4])
		copyVal, _ := strconv.Atoi(r[5])
		scanVal, _ := strconv.Atoi(r[6])

		dataMap[r[0]] = PrinterData{
			IP:       r[0],
			Location: r[1],
			Model:    r[2],
			Total:    total,
			Print:    printVal,
			Copy:     copyVal,
			Scan:     scanVal,
			Date:     r[10],
		}
	}
	return dataMap, nil
}

func main() {
	// 1. 取得所有 CSV 檔案
	dataPath := "private-data/data/*.csv" // 根據你的資料夾結構調整
	files, err := filepath.Glob(dataPath)
	if err != nil {
		log.Fatalf("搜尋路徑錯誤: %v", err)
	}

	// 2. 排序檔案 (按日期 yyyy-mm-dd.csv 排序)
	sort.Strings(files)

	// 檢查檔案數量
	if len(files) == 0 {
		fmt.Println("錯誤：找不到任何 CSV 數據檔案。")
		os.Exit(0) // 正常結束，避免 GitHub Action 顯示紅燈
	}

	if len(files) < 2 {
		fmt.Printf("提示：目前僅有 %d 個檔案，無法計算差額。需至少兩份數據。\n", len(files))
		// 如果只有一份檔案，我們可以直接輸出該份資料作為基準，或者結束
		os.Exit(0)
	}

	// 3. 讀取最後兩天的資料 (前一天 vs 今天)
	prevFile := files[len(files)-2]
	todayFile := files[len(files)-1]

	fmt.Printf("正在處理資料：[%s] 與 [%s]\n", prevFile, todayFile)

	prevDay, err := readCSV(prevFile)
	if err != nil {
		log.Printf("警告：讀取前一天資料失敗: %v", err)
		os.Exit(0)
	}

	today, err := readCSV(todayFile)
	if err != nil {
		log.Printf("警告：讀取今日資料失敗: %v", err)
		os.Exit(0)
	}

	// 4. 計算今日增量
	var dailyUsage []PrinterData
	for ip, now := range today {
		before, ok := prevDay[ip]
		if !ok {
			// 如果這台印表機是新加入的，沒有昨日數據，則增量設為 0 或今日總量
			dailyUsage = append(dailyUsage, PrinterData{
				IP:       now.IP,
				Location: now.Location,
				Model:    now.Model,
				Total:    0,
				Date:     now.Date,
			})
			continue
		}

		// 計算差值 (今日總數 - 昨日總數)
		// 注意：若遇到印表機計數器歸零（例如更換主機板），差值可能為負，這裡取 0
		calcUsage := func(n, b int) int {
			res := n - b
			if res < 0 {
				return 0
			}
			return res
		}

		dailyUsage = append(dailyUsage, PrinterData{
			IP:       now.IP,
			Location: now.Location,
			Model:    now.Model,
			Total:    calcUsage(now.Total, before.Total),
			Print:    calcUsage(now.Print, before.Print),
			Copy:     calcUsage(now.Copy, before.Copy),
			Scan:     calcUsage(now.Scan, before.Scan),
			Date:     now.Date,
		})
	}

	// 5. 輸出 JSON
	jsonData, err := json.MarshalIndent(dailyUsage, "", "  ")
	if err != nil {
		log.Fatalf("JSON 編碼失敗: %v", err)
	}

	err = os.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("寫入檔案失敗: %v", err)
	}

	fmt.Println("成功：data.json 已更新。")
}
