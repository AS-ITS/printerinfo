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

type PrinterStats struct {
	IP       string  `json:"ip"`
	Location string  `json:"location"`
	Model    string  `json:"model"`
	Daily    Metrics `json:"daily"`
	Monthly  Metrics `json:"monthly"`
	Yearly   Metrics `json:"yearly"`
	RawTotal int     `json:"raw_total"` // 用於內部計算
}

func readCSV(path string) (map[string]map[string]int, error) {
	data := make(map[string]map[string]int)
	if path == "" {
		return data, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	for i, r := range records {
		if i == 0 || len(r) < 11 {
			continue
		}
		ip := r[0]
		data[ip] = map[string]int{
			"total": r_int(r[3]),
			"print": r_int(r[4]),
			"copy":  r_int(r[5]),
			"fax":   r_int(r[9]),
		}
	}
	return data, nil
}

func r_int(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func calculateDelta(current, base map[string]int) Metrics {
	return Metrics{
		Print: current["print"] - base["print"],
		Copy:  current["copy"] - base["copy"],
		Fax:   current["fax"] - base["fax"],
		Total: current["total"] - base["total"],
	}
}

func main() {
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)
	if len(files) == 0 {
		return
	}

	latestFile := files[len(files)-1]
	currentMonth := filepath.Base(latestFile)[:7]
	currentYear := filepath.Base(latestFile)[:4]

	// 找出基準檔案
	var dBaseFile, mBaseFile, yBaseFile string
	if len(files) >= 2 {
		dBaseFile = files[len(files)-2]
	} else {
		dBaseFile = latestFile
	}

	for _, f := range files {
		if mBaseFile == "" && strings.HasPrefix(filepath.Base(f), currentMonth) {
			mBaseFile = f
		}
		if yBaseFile == "" && strings.HasPrefix(filepath.Base(f), currentYear) {
			yBaseFile = f
		}
	}

	latestData, _ := readCSV(latestFile)
	dBase, _ := readCSV(dBaseFile)
	mBase, _ := readCSV(mBaseFile)
	yBase, _ := readCSV(yBaseFile)

	// 取得基本資訊（從最新檔案讀取地點與型號）
	f, _ := os.Open(latestFile)
	records, _ := csv.NewReader(f).ReadAll()
	f.Close()

	var results []PrinterStats
	for i, r := range records {
		if i == 0 {
			continue
		}
		ip := r[0]
		stats := PrinterStats{
			IP:       ip,
			Location: r[1],
			Model:    r[2],
			Daily:    calculateDelta(latestData[ip], dBase[ip]),
			Monthly:  calculateDelta(latestData[ip], mBase[ip]),
			Yearly:   calculateDelta(latestData[ip], yBase[ip]),
			RawTotal: latestData[ip]["total"],
		}
		results = append(results, stats)
	}

	output, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("data.json", output, 0644)
	fmt.Println("data.json generated successfully.")
}
