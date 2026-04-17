package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type PrinterData struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
	Model    string `json:"model"`
	Total    int    `json:"total"`
	Print    int    `json:"print"`
	Copy     int    `json:"copy"`
	Scan     int    `json:"scan"`
	Date     string `json:"date"`
}

func readCSV(path string) (map[string]PrinterData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	dataMap := make(map[string]PrinterData)

	for i, r := range records {
		if i == 0 {
			continue
		} // 跳過標題
		total, _ := strconv.Atoi(r[3])
		printVal, _ := strconv.Atoi(r[4])
		copyVal, _ := strconv.Atoi(r[5])
		scanVal, _ := strconv.Atoi(r[6])

		dataMap[r[0]] = PrinterData{
			IP: r[0], Location: r[1], Model: r[2],
			Total: total, Print: printVal, Copy: copyVal, Scan: scanVal,
			Date: r[10],
		}
	}
	return dataMap, nil
}

func main() {
	files, _ := filepath.Glob("private-data/data/*.csv")
	sort.Strings(files)

	if len(files) < 2 {
		fmt.Println("數據不足，至少需要兩天的檔案")
		return
	}

	// 取得最後兩天的檔案
	prevDay, _ := readCSV(files[len(files)-2])
	today, _ := readCSV(files[len(files)-1])

	var dailyUsage []PrinterData

	for ip, now := range today {
		if before, ok := prevDay[ip]; ok {
			dailyUsage = append(dailyUsage, PrinterData{
				IP:       now.IP,
				Location: now.Location,
				Model:    now.Model,
				Total:    now.Total - before.Total,
				Print:    now.Print - before.Print,
				Copy:     now.Copy - before.Copy,
				Scan:     now.Scan - before.Scan,
				Date:     now.Date,
			})
		}
	}

	output, _ := json.MarshalIndent(dailyUsage, "", "  ")
	_ = ioutil.WriteFile("data.json", output, 0644)
	fmt.Println("JSON 轉換完成")
}
