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

type MonthData struct {
	Month   string  `json:"month"` // "01", "02" ...
	Metrics Metrics `json:"metrics"`
}

type YearData struct {
	Year    string      `json:"year"`
	Metrics Metrics     `json:"year_total"`
	Months  []MonthData `json:"months"`
}

type PrinterHistory struct {
	IP       string     `json:"ip"`
	Location string     `json:"location"`
	Model    string     `json:"model"`
	History  []YearData `json:"history"`
}

func r_int(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	// 1. 取得所有檔案並排序
	files, _ := filepath.Glob("*.csv") // 根據環境調整路徑，例如 "data/*.csv"
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("No CSV files found.")
		return
	}

	// 用於儲存每台印表機「前一個檔案」的數值，以便計算增量
	prevCounters := make(map[string]map[string]int)
	// 最終統計結果：IP -> Year -> Month -> Metrics
	stats := make(map[string]map[string]map[string]*Metrics)
	// 基本資訊儲存
	info := make(map[string][2]string) // IP -> [Location, Model]

	for _, file := range files {
		f, _ := os.Open(file)
		records, _ := csv.NewReader(f).ReadAll()
		f.Close()

		dateParts := strings.Split(filepath.Base(file), "-") // yyyy-mm-dd
		year := dateParts[0]
		month := dateParts[1]

		for i, r := range records {
			if i == 0 || len(r) < 11 {
				continue
			}
			ip := r[0]
			info[ip] = [2]string{r[1], r[2]}

			currTotal := r_int(r[3])
			currPrint := r_int(r[4])
			currCopy := r_int(r[5])
			currFax := r_int(r[9])

			// 如果有前一次紀錄，計算增量
			if prev, ok := prevCounters[ip]; ok {
				dTotal := currTotal - prev["total"]
				dPrint := currPrint - prev["print"]
				dCopy := currCopy - prev["copy"]
				dFax := currFax - prev["fax"]

				// 確保增量不為負數 (防止換碳粉匣或維修後計數器歸零的情況)
				if dTotal < 0 {
					dTotal = 0
				}
				if dPrint < 0 {
					dPrint = 0
				}
				if dCopy < 0 {
					dCopy = 0
				}
				if dFax < 0 {
					dFax = 0
				}

				// 初始化結構並累加
				if stats[ip] == nil {
					stats[ip] = make(map[string]map[string]*Metrics)
				}
				if stats[ip][year] == nil {
					stats[ip][year] = make(map[string]*Metrics)
				}
				if stats[ip][year][month] == nil {
					stats[ip][year][month] = &Metrics{}
				}

				stats[ip][year][month].Total += dTotal
				stats[ip][year][month].Print += dPrint
				stats[ip][year][month].Copy += dCopy
				stats[ip][year][month].Fax += dFax
			}

			// 更新前一次紀錄
			prevCounters[ip] = map[string]int{
				"total": currTotal, "print": currPrint, "copy": currCopy, "fax": currFax,
			}
		}
	}

	// 2. 轉換為前端易用的 JSON 結構
	var finalOutput []PrinterHistory
	for ip, years := range stats {
		pHistory := PrinterHistory{IP: ip, Location: info[ip][0], Model: info[ip][1]}

		var yearKeys []string
		for y := range years {
			yearKeys = append(yearKeys, y)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(yearKeys)))

		for _, y := range yearKeys {
			yData := YearData{Year: y}
			var monthKeys []string
			for m := range years[y] {
				monthKeys = append(monthKeys, m)
			}
			sort.Strings(monthKeys)

			for _, m := range monthKeys {
				mMetrics := *years[y][m]
				yData.Months = append(yData.Months, MonthData{Month: m, Metrics: mMetrics})
				// 累加年度總計
				yData.Metrics.Total += mMetrics.Total
				yData.Metrics.Print += mMetrics.Print
				yData.Metrics.Copy += mMetrics.Copy
				yData.Metrics.Fax += mMetrics.Fax
			}
			pHistory.History = append(pHistory.History, yData)
		}
		finalOutput = append(finalOutput, pHistory)
	}

	out, _ := json.MarshalIndent(finalOutput, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("History data.json generated successfully.")
}
