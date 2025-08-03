package main

import (
	"fmt"
)

func main() {
	fileTop := LoadFromCSV("./iranian_domains_from_top10m.csv")
	fileIran := LoadFromXLS("./iran.XLSX")
	PrintCommonElements(fileTop, fileIran)
}

func PrintCommonElements(fileTop, fileIran []string) {
	set := make(map[string]struct{})
	for _, domain := range fileTop {
		set[domain] = struct{}{}
	}

	for _, domain := range fileIran {
		fmt.Println(domain)
		if _, exists := set[domain]; exists {
			fmt.Println(domain)
		}
	}
}
