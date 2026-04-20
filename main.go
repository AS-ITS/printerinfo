package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Metrics struct {
	Print int `json:"print"`
	Copy  int `json:"copy"`
	Fax   int `json:"fax"`
	Total int `json:"total"`
}

type DailyLog struct {
	Date    string  `json:"date"`
	Metrics Metrics `json:"metrics"`
}

type MonthData struct {
	Month        string     `json:"month"`
	MonthMetrics Metrics    `json:"month_metrics"`
	DailyLogs    []DailyLog `json:"daily_logs"`
}

type YearData struct {
	Year      string      `json:"year"`
	YearTotal Metrics     `json:"year_total"`
	Months    []MonthData `json:"months"`
}

type PrinterStats struct {
	IP       string     `json:"ip"`
	Location string     `json:"location"`
	Model    string     `json:"model"`
	History  []YearData `json:"history"`
	Today    Metrics    `json:"today"`
}

type DashboardData struct {
	GeneratedAt string         `json:"generated_at"`
	Printers    []PrinterStats `json:"printers"`
}

func r_int(s string) int {
	// 移除千分號與空白
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	// 1. 取得所有檔案並嚴格排序
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("Error: No CSV files found in data/ folder.")
		return
	}

	fmt.Printf("Processing %d files: %v\n", len(files), files)

	// 用於追蹤前一個狀態
	prevCounters := make(map[string]map[string]int)
	// IP -> Year -> Month -> *MonthData
	storage := make(map[string]map[string]map[string]*MonthData)
	info := make(map[string][2]string)
	todayMetrics := make(map[string]Metrics)
	var latestDate string

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		records, _ := csv.NewReader(f).ReadAll()
		f.Close()

		base := filepath.Base(file)
		if len(base) < 10 {
			continue
		}
		dateStr := base[:10]
		latestDate = dateStr
		year, month := dateStr[:4], dateStr[5:7]

		for i, r := range records {
			if i == 0 || len(r) < 11 {
				continue
			}

			// 移除 IP 可能存在的隱形空白
			ip := strings.TrimSpace(r[0])
			if ip == "" {
				continue
			}

			info[ip] = [2]string{strings.TrimSpace(r[1]), strings.TrimSpace(r[2])}

			curr := map[string]int{
				"total": r_int(r[3]),
				"print": r_int(r[4]),
				"copy":  r_int(r[5]),
				"fax":   r_int(r[9]),
			}

			// 如果有前一個檔案的數據，計算增量
			if prev, ok := prevCounters[ip]; ok {
				dm := Metrics{
					Print: curr["print"] - prev["print"],
					Copy:  curr["copy"] - prev["copy"],
					Fax:   curr["fax"] - prev["fax"],
					Total: curr["total"] - prev["total"],
				}

				// 正常情況下計數器只會增加。若為負值則代表計數器重置或異常，設為 0
				if dm.Total < 0 {
					dm = Metrics{}
				}

				// 初始化結構
				if storage[ip] == nil {
					storage[ip] = make(map[string]map[string]*MonthData)
				}
				if storage[ip][year] == nil {
					storage[ip][year] = make(map[string]*MonthData)
				}
				if storage[ip][year][month] == nil {
					storage[ip][year][month] = &MonthData{Month: month, DailyLogs: []DailyLog{}}
				}

				m := storage[ip][year][month]
				// 檢查是否已存在同日期的 Log (避免重複處理)
				exists := false
				for _, log := range m.DailyLogs {
					if log.Date == dateStr {
						exists = true
						break
					}
				}

				if !exists {
					m.DailyLogs = append(m.DailyLogs, DailyLog{Date: dateStr, Metrics: dm})
					m.MonthMetrics.Print += dm.Print
					m.MonthMetrics.Copy += dm.Copy
					m.MonthMetrics.Fax += dm.Fax
					m.MonthMetrics.Total += dm.Total
					// 更新該設備最後一次的增量 (用於主頁表格)
					todayMetrics[ip] = dm
				}
			}
			// 更新基準值
			prevCounters[ip] = curr
		}
	}

	// 2. 組裝最終 JSON
	var printerList []PrinterStats
	for ip, years := range storage {
		p := PrinterStats{
			IP: ip, Location: info[ip][0], Model: info[ip][1],
			Today: todayMetrics[ip], History: []YearData{},
		}

		yKeys := make([]string, 0, len(years))
		for y := range years {
			yKeys = append(yKeys, y)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(yKeys)))

		for _, y := range yKeys {
			yData := YearData{Year: y, Months: []MonthData{}}
			mKeys := make([]string, 0, len(years[y]))
			for m := range years[y] {
				mKeys = append(mKeys, m)
			}
			sort.Strings(mKeys)

			for _, mk := range mKeys {
				mDat := years[y][mk]
				// 輸出前再次確保日期排序正確
				sort.Slice(mDat.DailyLogs, func(i, j int) bool {
					return mDat.DailyLogs[i].Date < mDat.DailyLogs[j].Date
				})
				yData.Months = append(yData.Months, *mDat)

				yData.YearTotal.Print += mDat.MonthMetrics.Print
				yData.YearTotal.Copy += mDat.MonthMetrics.Copy
				yData.YearTotal.Fax += mDat.MonthMetrics.Fax
				yData.YearTotal.Total += mDat.MonthMetrics.Total
			}
			p.History = append(p.History, yData)
		}
		printerList = append(printerList, p)
	}

	finalData := DashboardData{GeneratedAt: latestDate, Printers: printerList}
	out, _ := json.MarshalIndent(finalData, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("Success: data.json generated with all historical deltas.")
}
