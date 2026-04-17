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

type PrinterStats struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
	Total    int    `json:"total"` // 這裡指增量
	Date     string `json:"date"`
}

type OutputData struct {
	Daily   []PrinterStats `json:"daily"`
	Monthly []PrinterStats `json:"monthly"`
	Yearly  []PrinterStats `json:"yearly"`
}

func readCSV(path string) (map[string]PrinterStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	data := make(map[string]PrinterStats)
	for i, r := range records {
		if i == 0 || len(r) < 11 {
			continue
		}
		total, _ := strconv.Atoi(r[3])
		data[r[0]] = PrinterStats{IP: r[0], Location: r[1], Total: total, Date: r[10]}
	}
	return data, nil
}

// 計算差值並排序
func calculateDiff(latest, base map[string]PrinterStats) []PrinterStats {
	var result []PrinterStats
	for ip, now := range latest {
		val := now.Total
		if prev, ok := base[ip]; ok {
			val = now.Total - prev.Total
		}
		if val < 0 {
			val = 0
		}
		result = append(result, PrinterStats{IP: ip, Location: now.Location, Total: val, Date: now.Date})
	}
	// 依照 Total 排序 (由大到小)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})
	return result
}

func main() {
	files, _ := filepath.Glob("private-data/*.csv")
	sort.Strings(files)
	if len(files) < 1 {
		return
	}

	latestFile := files[len(files)-1]
	latestData, _ := readCSV(latestFile)

	// 找到基準檔案
	var prevDayFile, monthStartFile, yearStartFile string
	currentDate := filepath.Base(latestFile) // yyyy-mm-dd.csv
	currentMonth := currentDate[:7]          // yyyy-mm
	currentYear := currentDate[:4]           // yyyy

	for _, f := range files {
		name := filepath.Base(f)
		if name < currentDate && (prevDayFile == "" || name > filepath.Base(prevDayFile)) {
			prevDayFile = f
		}
		if strings.HasPrefix(name, currentMonth) && (monthStartFile == "" || name < filepath.Base(monthStartFile)) {
			monthStartFile = f
		}
		if strings.HasPrefix(name, currentYear) && (yearStartFile == "" || name < filepath.Base(yearStartFile)) {
			yearStartFile = f
		}
	}

	// 讀取基準數據
	dBase, _ := readCSV(prevDayFile)
	mBase, _ := readCSV(monthStartFile)
	yBase, _ := readCSV(yearStartFile)

	output := OutputData{
		Daily:   calculateDiff(latestData, dBase),
		Monthly: calculateDiff(latestData, mBase),
		Yearly:  calculateDiff(latestData, yBase),
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	_ = os.WriteFile("data.json", jsonBytes, 0644)
	fmt.Println("統計完成，已依照列印量排序。")
}
