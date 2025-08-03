package main

import (
	"encoding/csv"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/publicsuffix"
	"io"
	"net"
	"os"
	"strings"
)

func LoadFromCSV(path string) []string {
	unique := make(map[string]struct{})
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			// ignore error from Fprintf to os.Stderr; non-critical
			_, _ = fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
		}
	}()

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

		if len(row) < 2 || row[1] == "" {
			continue
		}

		if _, exists := unique[row[1]]; !exists {
			unique[row[1]] = struct{}{}
			domains = append(domains, row[1])
		}
	}
	return domains
}

func LoadFromXLS(path string) []string {
	f, err := excelize.OpenFile(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var domains []string
	unique := make(map[string]struct{})
	rows, err := f.GetRows("iran")

	if err != nil {
		panic(err)
	}
	for index, row := range rows {
		if index == 0 || len(row) < 2 || row[1] == "" {
			continue
		}

		cleanRow := strings.TrimSpace(row[0])
		cleanRow = strings.TrimPrefix(cleanRow, "http://")
		cleanRow = strings.TrimPrefix(cleanRow, "https://")
		cleanRow = strings.TrimRight(cleanRow, ".,/")

		if cleanRow == "" {
			continue
		}

		// to ignore ip like domains like http://185.112.33.51
		if net.ParseIP(cleanRow) != nil {
			continue
		}

		sld, err := publicsuffix.EffectiveTLDPlusOne(cleanRow)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if _, exist := unique[sld]; !exist {
			unique[sld] = struct{}{}
			domains = append(domains, sld)
		}
	}
	return domains
}
