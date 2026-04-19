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

func readCSV(path string) (map[string]map[string]int, error) {
	data := make(map[string]map[string]int)
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

func r_int(s string) int { v, _ := strconv.Atoi(s); return v }

func calculateDelta(cur, base map[string]int) Metrics {
	if cur == nil || base == nil {
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
		return
	}

	latestFile := files[len(files)-1]
	currentMonth := filepath.Base(latestFile)[:7]
	currentYear := filepath.Base(latestFile)[:4]

	var dBaseFile, mBaseFile, yBaseFile string
	if len(files) >= 2 {
		dBaseFile = files[len(files)-2]
	} else {
		dBaseFile = latestFile
	}
	for _, f := range files {
		base := filepath.Base(f)
		if mBaseFile == "" && strings.HasPrefix(base, currentMonth) {
			mBaseFile = f
		}
		if yBaseFile == "" && strings.HasPrefix(base, currentYear) {
			yBaseFile = f
		}
	}

	latest, _ := readCSV(latestFile)
	dBase, _ := readCSV(dBaseFile)
	mBase, _ := readCSV(mBaseFile)
	yBase, _ := readCSV(yBaseFile)

	f, _ := os.Open(latestFile)
	records, _ := csv.NewReader(f).ReadAll()
	f.Close()

	var results []PrinterStats
	for i, r := range records {
		if i == 0 || len(r) < 11 {
			continue
		}
		ip := r[0]
		stats := PrinterStats{
			IP: ip, Location: r[1], Model: r[2],
			Daily:    calculateDelta(latest[ip], dBase[ip]),
			Monthly:  calculateDelta(latest[ip], mBase[ip]),
			Yearly:   calculateDelta(latest[ip], yBase[ip]),
			RawTotal: latest[ip]["total"],
		}
		results = append(results, stats)
	}

	output, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("data.json", output, 0644)
	fmt.Println("data.json generated.")
}
