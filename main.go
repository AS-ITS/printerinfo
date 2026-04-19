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
	Year   string      `json:"year"`
	Months []MonthData `json:"months"`
}

type PrinterStats struct {
	IP       string     `json:"ip"`
	Location string     `json:"location"`
	Model    string     `json:"model"`
	History  []YearData `json:"history"`
	Today    Metrics    `json:"today"` // 當日最新增量
}

func r_int(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	files, _ := filepath.Glob("data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("Error: No CSV files in data/")
		return
	}

	prevCounters := make(map[string]map[string]int)
	// IP -> Year -> Month -> MonthData
	storage := make(map[string]map[string]map[string]*MonthData)
	info := make(map[string][2]string)
	todayMetrics := make(map[string]Metrics)

	for _, file := range files {
		f, _ := os.Open(file)
		records, _ := csv.NewReader(f).ReadAll()
		f.Close()

		base := filepath.Base(file) // yyyy-mm-dd.csv
		dateStr := base[:10]
		year, month := dateStr[:4], dateStr[5:7]

		for i, r := range records {
			if i == 0 || len(r) < 11 {
				continue
			}
			ip := r[0]
			info[ip] = [2]string{r[1], r[2]}

			curr := map[string]int{
				"total": r_int(r[3]), "print": r_int(r[4]), "copy": r_int(r[5]), "fax": r_int(r[9]),
			}

			if prev, ok := prevCounters[ip]; ok {
				dm := Metrics{
					Print: curr["print"] - prev["print"],
					Copy:  curr["copy"] - prev["copy"],
					Fax:   curr["fax"] - prev["fax"],
					Total: curr["total"] - prev["total"],
				}
				if dm.Total < 0 {
					dm = Metrics{}
				} // 排除異常

				// 儲存至月份與每日紀錄
				if storage[ip] == nil {
					storage[ip] = make(map[string]map[string]*MonthData)
				}
				if storage[ip][year] == nil {
					storage[ip][year] = make(map[string]*Metrics)
				} // 修正結構

				// 確保月份結構存在
				if _, ok := storage[ip][year][month]; !ok {
					storage[ip][year][month] = &MonthData{Month: month}
				}

				m := storage[ip][year][month]
				m.DailyLogs = append(m.DailyLogs, DailyLog{Date: dateStr, Metrics: dm})
				m.MonthMetrics.Print += dm.Print
				m.MonthMetrics.Copy += dm.Copy
				m.MonthMetrics.Fax += dm.Fax
				m.MonthMetrics.Total += dm.Total

				// 紀錄為最後一天的「今日量」
				todayMetrics[ip] = dm
			}
			prevCounters[ip] = curr
		}
	}

	// 轉化為前端格式
	var result []PrinterStats
	for ip, years := range storage {
		p := PrinterStats{IP: ip, Location: info[ip][0], Model: info[ip][1], Today: todayMetrics[ip]}
		for y, months := range years {
			yData := YearData{Year: y}
			var mKeys []string
			for m := range months {
				mKeys = append(mKeys, m)
			}
			sort.Strings(mKeys)
			for _, mk := range mKeys {
				yData.Months = append(yData.Months, *months[mk])
			}
			p.History = append(p.History, yData)
		}
		result = append(result, p)
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("data.json generated with daily details.")
}
