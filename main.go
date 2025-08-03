package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang/v2"
)

func main() {
	db, err := geoip2.Open("./GeoLite2-City.mmdb")
	defer func(db *geoip2.Reader) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	if err != nil {
		panic(err)
	}
	// TODO: Return Errors for all functions and main should handle the errors
	domains10M := LoadFromCSV("./list10m.csv")
	//fmt.Println(domains10M[432])
	_ = LoadFromXLS("./iran.xlsx")

	writeChannel := make(chan []string)
	wg := sync.WaitGroup{}

	// Saving irDomains as csv in a file
	outputFile, err := os.Create("iranian_domains_from_top10m.csv")
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(outputFile)

	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(outputFile)
	writer.Write([]string{"Domain"})
	defer writer.Flush()
	// Here we make 10_000 goroutine
	chunkSize := 1000
	for i := 0; i < len(domains10M); i += chunkSize {
		end := i + chunkSize
		if end > len(domains10M) {
			end = len(domains10M)
		}
		chunk := domains10M[i:end]
		wg.Add(1)
		go func(domains []string) {
			defer wg.Done()
			irDomains := FindIranFromTop10M(domains, db)
			if len(irDomains) > 0 {
				writeChannel <- irDomains
			}
		}(chunk)
	}
	go func() {
		wg.Wait()
		close(writeChannel)
	}()

	for domains := range writeChannel {
		for _, domain := range domains {
			if err := writer.Write([]string{domain}); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing record: %v\n", err)
			}
		}
	}
}

// TODO: open db one time and import it to each goroutine
