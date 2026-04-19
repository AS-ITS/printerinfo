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
	Month   string  `json:"month"`
	Metrics Metrics `json:"metrics"`
}

type YearData struct {
	Year      string      `json:"year"`
	YearTotal Metrics     `json:"year_total"`
	Months    []MonthData `json:"months"`
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
	// 假設 CSV 都在同一個目錄
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("Error: No CSV files found")
		return
	}

	// 紀錄每台印表機在前一個檔案的讀數
	prevCounters := make(map[string]map[string]int)
	// IP -> Year -> Month -> Metrics
	rawStats := make(map[string]map[string]map[string]*Metrics)
	info := make(map[string][2]string) // IP -> [Location, Model]

	for _, file := range files {
		f, _ := os.Open(file)
		records, _ := csv.NewReader(f).ReadAll()
		f.Close()

		// 檔名格式: yyyy-mm-dd.csv
		baseName := filepath.Base(file)
		if len(baseName) < 10 {
			continue
		}
		year := baseName[:4]
		month := baseName[5:7]

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

			if prev, ok := prevCounters[ip]; ok {
				// 計算今日增量
				dTotal := curr["total"] - prev["total"]
				dPrint := curr["print"] - prev["print"]
				dCopy := curr["copy"] - prev["copy"]
				dFax := curr["fax"] - prev["fax"]

				// 排除計數器歸零或異常負值
				if dTotal < 0 {
					dTotal = 0
				}

				// 累加到該月份
				if rawStats[ip] == nil {
					rawStats[ip] = make(map[string]map[string]*Metrics)
				}
				if rawStats[ip][year] == nil {
					rawStats[ip][year] = make(map[string]*Metrics)
				}
				if rawStats[ip][year][month] == nil {
					rawStats[ip][year][month] = &Metrics{}
				}

				rawStats[ip][year][month].Total += dTotal
				rawStats[ip][year][month].Print += dPrint
				rawStats[ip][year][month].Copy += dCopy
				rawStats[ip][year][month].Fax += dFax
			}
			prevCounters[ip] = curr
		}
	}

	// 轉換成前端用的結構
	var result []PrinterHistory
	for ip, years := range rawStats {
		p := PrinterHistory{IP: ip, Location: info[ip][0], Model: info[ip][1]}

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
				mMet := *years[y][m]
				yData.Months = append(yData.Months, MonthData{Month: m, Metrics: mMet})
				yData.YearTotal.Total += mMet.Total
				yData.YearTotal.Print += mMet.Print
				yData.YearTotal.Copy += mMet.Copy
				yData.YearTotal.Fax += mMet.Fax
			}
			p.History = append(p.History, yData)
		}
		result = append(result, p)
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("Done: data.json generated.")
}
