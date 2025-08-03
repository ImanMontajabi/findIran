package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fileTop := LoadFromCSVSingleColumn("./iranian_domains_from_top10m.csv")
	fileIranXLS := LoadFromXLS("./iran.xlsx")
	Save_xlsv_as_csv(fileIranXLS)
	fileIranCSV := LoadFromCSVSingleColumn("./iranian_domains.csv")
	fmt.Printf("Loaded %d domains from top10m file\n", len(fileTop))
	fmt.Printf("Loaded %d domains from Excel file\n", len(fileIranXLS))
	fmt.Printf("Loaded %d domains from CSV file\n", len(fileIranCSV))
	PrintDifferentElements(fileTop, fileIranCSV)
	err := SaveDifferentElementsToCSV(fileTop, fileIranCSV, "different_domains.csv")
	if err != nil {
		panic(err)
	}
}

func LoadFromCSVSingleColumn(path string) []string {
	unique := make(map[string]struct{})
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	// skip header
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	var domains []string

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if len(row) < 1 || row[0] == "" {
			continue
		}

		if _, exists := unique[row[0]]; !exists {
			unique[row[0]] = struct{}{}
			domains = append(domains, row[0])
		}
	}
	return domains
}

func PrintDifferentElements(fileTop, fileIran []string) {
	set := make(map[string]struct{})
	for _, domain := range fileIran {
		set[domain] = struct{}{}
	}

	for _, domain := range fileTop {
		if _, exists := set[domain]; !exists {
			fmt.Println(domain)
		}
	}
}

func Save_xlsv_as_csv(fileIran []string) {
	outputFile, err := os.Create("iranian_domains.csv")
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(outputFile)
	writer.Write([]string{"Domain"})
	for _, domain := range fileIran {
		writer.Write([]string{domain})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}
}

func SaveDifferentElementsToCSV(fileTop, fileIran []string, outPath string) error {
	set := make(map[string]struct{})
	for _, domain := range fileIran {
		d := strings.TrimSpace(strings.ToLower(domain))
		set[d] = struct{}{}
	}

	outputFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	writer.Write([]string{"Domain"}) // optional header

	for _, domain := range fileTop {
		d := strings.TrimSpace(strings.ToLower(domain))
		if _, exists := set[d]; exists {
			err := writer.Write([]string{d})
			if err != nil {
				return err
			}
		}
	}

	return writer.Error()
}
