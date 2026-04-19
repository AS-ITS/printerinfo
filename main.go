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
	RawTotal int     `json:"raw_total"`
}

func readData(path string) (map[string]map[string]int, error) {
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
	v, _ := strconv.Atoi(strings.ReplaceAll(s, ",", ""))
	return v
}

func calcDelta(cur, base map[string]int) Metrics {
	if cur == nil {
		return Metrics{}
	}
	if base == nil {
		return Metrics{}
	}
	return Metrics{
		Print: cur["print"] - base["print"],
		Copy:  cur["copy"] - base["copy"],
		Fax:   cur["fax"] - base["fax"],
		Total: cur["total"] - base["total"],
	}
}

func main() {
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)
	if len(files) == 0 {
		fmt.Println("No CSV files found")
		return
	}

	lastIdx := len(files) - 1
	latestFile := files[lastIdx]
	curYm := filepath.Base(latestFile)[:7] // yyyy-mm
	curY := filepath.Base(latestFile)[:4]  // yyyy

	var dBase, mBase, yBase string

	// 1. 日基準：前一個檔案
	if lastIdx > 0 {
		dBase = files[lastIdx-1]
	} else {
		dBase = latestFile
	}

	// 2. 月基準：找「非本月」的最後一個檔案
	for i := lastIdx; i >= 0; i-- {
		if !strings.HasPrefix(filepath.Base(files[i]), curYm) {
			mBase = files[i]
			break
		}
	}
	// 如果沒找到上個月的，就找本月第一個
	if mBase == "" {
		for _, f := range files {
			if strings.HasPrefix(filepath.Base(f), curYm) {
				mBase = f
				break
			}
		}
	}

	// 3. 年基準：找「非本年」的最後一個檔案
	for i := lastIdx; i >= 0; i-- {
		if !strings.HasPrefix(filepath.Base(files[i]), curY) {
			yBase = files[i]
			break
		}
	}
	if yBase == "" {
		for _, f := range files {
			if strings.HasPrefix(filepath.Base(f), curY) {
				yBase = f
				break
			}
		}
	}

	latestMap, _ := readData(latestFile)
	dBaseMap, _ := readData(dBase)
	mBaseMap, _ := readData(mBase)
	yBaseMap, _ := readData(yBase)

	f, _ := os.Open(latestFile)
	records, _ := csv.NewReader(f).ReadAll()
	f.Close()

	var results []PrinterStats
	for i, r := range records {
		if i == 0 || len(r) < 11 {
			continue
		}
		ip := r[0]
		cur := latestMap[ip]

		results = append(results, PrinterStats{
			IP: ip, Location: r[1], Model: r[2],
			Daily:    calcDelta(cur, dBaseMap[ip]),
			Monthly:  calcDelta(cur, mBaseMap[ip]),
			Yearly:   calcDelta(cur, yBaseMap[ip]),
			RawTotal: cur["total"],
		})
	}

	out, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("data.json", out, 0644)
	fmt.Println("Calculation complete. data.json saved.")
}
