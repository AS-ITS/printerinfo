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
	Model    string `json:"model"`
	Total    int    `json:"total"`
	Date     string `json:"date"`
}

type OutputData struct {
	Daily   []PrinterStats    `json:"daily"`
	Monthly []PrinterStats    `json:"monthly"`
	Yearly  []PrinterStats    `json:"yearly"`
	Meta    map[string]string `json:"meta"`
}

func readCSV(path string) (map[string]PrinterStats, error) {
	if path == "" {
		return make(map[string]PrinterStats), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	data := make(map[string]PrinterStats)
	for i, r := range records {
		// 跳過標題且確保欄位數正確 (根據範例 11 欄)
		if i == 0 || len(r) < 11 {
			continue
		}
		total, _ := strconv.Atoi(r[3])
		data[r[0]] = PrinterStats{
			IP:       r[0],
			Location: r[1],
			Model:    r[2],
			Total:    total,
			Date:     r[10],
		}
	}
	return data, nil
}

func calculateDiff(latest, base map[string]PrinterStats) []PrinterStats {
	var result []PrinterStats
	for ip, now := range latest {
		usage := 0
		if prev, ok := base[ip]; ok {
			usage = now.Total - prev.Total
		}
		if usage < 0 {
			usage = 0
		}
		result = append(result, PrinterStats{
			IP:       ip,
			Location: now.Location,
			Model:    now.Model,
			Total:    usage,
			Date:     now.Date,
		})
	}
	// 依照使用量從大到小排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})
	return result
}

func main() {
	// 假設 Private Repo 檢出在 private-data 目錄下
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Println("錯誤：找不到資料檔案")
		return
	}

	latestFile := files[len(files)-1]
	earliestFile := files[0]

	// 取得當前檔案的時間維度
	currentName := filepath.Base(latestFile)
	currentMonth := currentName[:7] // yyyy-mm
	currentYear := currentName[:4]  // yyyy

	// 1. 尋找日基準 (上一個可用的檔案，解決週末不蒐集的問題)
	var dBaseFile string
	if len(files) >= 2 {
		dBaseFile = files[len(files)-2]
	} else {
		dBaseFile = latestFile
	}

	// 2. 尋找月基準 (本月最早的一份)
	mBaseFile := earliestFile
	for _, f := range files {
		if strings.HasPrefix(filepath.Base(f), currentMonth) {
			mBaseFile = f
			break
		}
	}

	// 3. 尋找年基準 (本年最早的一份)
	yBaseFile := earliestFile
	for _, f := range files {
		if strings.HasPrefix(filepath.Base(f), currentYear) {
			yBaseFile = f
			break
		}
	}

	latestData, _ := readCSV(latestFile)
	dBase, _ := readCSV(dBaseFile)
	mBase, _ := readCSV(mBaseFile)
	yBase, _ := readCSV(yBaseFile)

	output := OutputData{
		Daily:   calculateDiff(latestData, dBase),
		Monthly: calculateDiff(latestData, mBase),
		Yearly:  calculateDiff(latestData, yBase),
		Meta: map[string]string{
			"latest":  filepath.Base(latestFile),
			"daily":   filepath.Base(dBaseFile),
			"monthly": filepath.Base(mBaseFile),
			"yearly":  filepath.Base(yBaseFile),
		},
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	_ = os.WriteFile("data.json", jsonBytes, 0644)
	fmt.Println("資料預處理完成：data.json 已更新")
}
