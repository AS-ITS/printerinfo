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
	Today    Metrics    `json:"today"`
}

func r_int(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("Error: No CSV files in data/")
		return
	}

	prevCounters := make(map[string]map[string]int)
	// IP -> Year -> Month -> *MonthData
	storage := make(map[string]map[string]map[string]*MonthData)
	info := make(map[string][2]string)
	todayMetrics := make(map[string]Metrics)

	for _, file := range files {
		f, _ := os.Open(file)
		records, _ := csv.NewReader(f).ReadAll()
		f.Close()

		base := filepath.Base(file) // yyyy-mm-dd.csv
		if len(base) < 10 {
			continue
		}
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
				}

				// 初始化 IP 層
				if storage[ip] == nil {
					storage[ip] = make(map[string]map[string]*MonthData)
				}
				// 初始化年份層
				if storage[ip][year] == nil {
					storage[ip][year] = make(map[string]*MonthData) // 修正點：必須與結構定義一致
				}
				// 初始化月份層
				if storage[ip][year][month] == nil {
					storage[ip][year][month] = &MonthData{Month: month, DailyLogs: []DailyLog{}}
				}

				m := storage[ip][year][month]
				m.DailyLogs = append(m.DailyLogs, DailyLog{Date: dateStr, Metrics: dm})
				m.MonthMetrics.Print += dm.Print
				m.MonthMetrics.Copy += dm.Copy
				m.MonthMetrics.Fax += dm.Fax
				m.MonthMetrics.Total += dm.Total

				todayMetrics[ip] = dm
			}
			prevCounters[ip] = curr
		}
	}

	var result []PrinterStats
	for ip, years := range storage {
		p := PrinterStats{IP: ip, Location: info[ip][0], Model: info[ip][1], Today: todayMetrics[ip], History: []YearData{}}

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
				yData.Months = append(yData.Months, *years[y][mk])
			}
			p.History = append(p.History, yData)
		}
		result = append(result, p)
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("Success: data.json generated.")
}
