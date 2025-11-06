package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Jsanchez767/InfluencePower/backend/cityapi"
)

func main() {
	client := cityapi.NewClient()
	
	records, err := client.GetOfficeRecords()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Total office records: %d\n\n", len(records))
	
	// Count by body name
	bodyCount := make(map[string]int)
	for _, record := range records {
		bodyCount[record.OfficeRecordBodyName]++
	}
	
	fmt.Println("Records by body:")
	for body, count := range bodyCount {
		fmt.Printf("  %s: %d\n", body, count)
	}
	
	// Find current City Council members
	now := time.Now()
	var currentCityCouncil []cityapi.OfficeRecord
	
	for _, record := range records {
		if record.OfficeRecordBodyName != "City Council" {
			continue
		}
		
		// Check if currently serving
		if record.OfficeRecordEndDate == "" {
			currentCityCouncil = append(currentCityCouncil, record)
			continue
		}
		
		endDate, err := time.Parse("2006-01-02T15:04:05", record.OfficeRecordEndDate)
		if err == nil && endDate.After(now) {
			currentCityCouncil = append(currentCityCouncil, record)
		}
	}
	
	fmt.Printf("\nCurrent City Council members: %d\n\n", len(currentCityCouncil))
	
	if len(currentCityCouncil) > 0 {
		fmt.Println("First 5 current officials:")
		for i, record := range currentCityCouncil {
			if i >= 5 {
				break
			}
			fmt.Printf("  - %s (%s)\n", record.OfficeRecordFullName, record.OfficeRecordTitle)
			fmt.Printf("    Email: %s\n", record.OfficeRecordEmail)
			fmt.Printf("    Start: %s, End: %s\n", record.OfficeRecordStartDate, record.OfficeRecordEndDate)
			fmt.Println()
		}
	}
	
	// Show sample of all records
	fmt.Println("Sample of all records (first 3):")
	for i, record := range records {
		if i >= 3 {
			break
		}
		fmt.Printf("  - %s (%s) - %s\n", record.OfficeRecordFullName, record.OfficeRecordTitle, record.OfficeRecordBodyName)
		fmt.Printf("    Start: %s, End: %s\n", record.OfficeRecordStartDate, record.OfficeRecordEndDate)
		fmt.Println()
	}
}
