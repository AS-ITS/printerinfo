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
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	// 1. 取得所有 CSV 檔案並排序
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("Error: No CSV files found in data/ folder.")
		return
	}

	fmt.Printf("Found %d files. Processing...\n", len(files))

	// 用於追蹤前一個狀態
	prevCounters := make(map[string]map[string]int)
	// 核心儲存結構: IP -> Year -> Month -> *MonthData
	storage := make(map[string]map[string]map[string]*MonthData)
	info := make(map[string][2]string)
	todayMetrics := make(map[string]Metrics)
	var latestDate string

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Skip file %s: %v\n", file, err)
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
			ip := r[0]
			info[ip] = [2]string{r[1], r[2]}

			curr := map[string]int{
				"total": r_int(r[3]),
				"print": r_int(r[4]),
				"copy":  r_int(r[5]),
				"fax":   r_int(r[9]),
			}

			// 如果不是第一份檔案，就計算與上一個檔案的增量
			if prev, ok := prevCounters[ip]; ok {
				dm := Metrics{
					Print: curr["print"] - prev["print"],
					Copy:  curr["copy"] - prev["copy"],
					Fax:   curr["fax"] - prev["fax"],
					Total: curr["total"] - prev["total"],
				}
				// 排除負值（可能是計數器重置）
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

				// 記錄每日紀錄
				m := storage[ip][year][month]
				m.DailyLogs = append(m.DailyLogs, DailyLog{Date: dateStr, Metrics: dm})

				// 累加月份指標
				m.MonthMetrics.Print += dm.Print
				m.MonthMetrics.Copy += dm.Copy
				m.MonthMetrics.Fax += dm.Fax
				m.MonthMetrics.Total += dm.Total

				// 由於是按日期排序，最後更新的會是「今日增量」
				todayMetrics[ip] = dm
			}
			// 更新前一個計數器基準
			prevCounters[ip] = curr
		}
	}

	// 2. 轉換為輸出結構
	var printerList []PrinterStats
	for ip, years := range storage {
		p := PrinterStats{IP: ip, Location: info[ip][0], Model: info[ip][1], Today: todayMetrics[ip]}

		yKeys := make([]string, 0, len(years))
		for y := range years {
			yKeys = append(yKeys, y)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(yKeys)))

		for _, y := range yKeys {
			yData := YearData{Year: y}
			mKeys := make([]string, 0, len(years[y]))
			for m := range years[y] {
				mKeys = append(mKeys, m)
			}
			sort.Strings(mKeys)

			for _, mk := range mKeys {
				mDat := *years[y][mk]
				// 排序每日紀錄
				sort.Slice(mDat.DailyLogs, func(i, j int) bool {
					return mDat.DailyLogs[i].Date < mDat.DailyLogs[j].Date
				})

				yData.Months = append(yData.Months, mDat)
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
	fmt.Println("Processed all files and generated data.json successfully.")
}
