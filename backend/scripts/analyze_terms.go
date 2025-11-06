package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Jsanchez767/InfluencePower/backend/cityapi"
)

func main() {
	client := cityapi.NewClient()
	
	records, err := client.GetOfficeRecords()
	if err != nil {
		log.Fatal(err)
	}
	
	// Filter City Council only
	var cityCouncil []cityapi.OfficeRecord
	for _, record := range records {
		if record.OfficeRecordBodyName == "City Council" {
			cityCouncil = append(cityCouncil, record)
		}
	}
	
	fmt.Printf("Total City Council records: %d\n\n", len(cityCouncil))
	
	// Sort by start date (most recent first)
	sort.Slice(cityCouncil, func(i, j int) bool {
		return cityCouncil[i].OfficeRecordStartDate > cityCouncil[j].OfficeRecordStartDate
	})
	
	fmt.Println("Most recent City Council records:")
	for i, record := range cityCouncil {
		if i >= 15 {
			break
		}
		
		endDate := record.OfficeRecordEndDate
		if endDate == "" {
			endDate = "CURRENT"
		}
		
		fmt.Printf("%2d. %s (%s)\n", i+1, record.OfficeRecordFullName, record.OfficeRecordTitle)
		fmt.Printf("    Email: %s\n", record.OfficeRecordEmail)
		fmt.Printf("    Start: %s → End: %s\n", record.OfficeRecordStartDate, endDate)
		fmt.Println()
	}
	
	// Check for records with no end date
	var noEndDate []cityapi.OfficeRecord
	for _, record := range cityCouncil {
		if record.OfficeRecordEndDate == "" {
			noEndDate = append(noEndDate, record)
		}
	}
	
	fmt.Printf("Records with no end date: %d\n", len(noEndDate))
	
	// Check for records ending after 2023
	var recent []cityapi.OfficeRecord
	for _, record := range cityCouncil {
		if record.OfficeRecordEndDate == "" {
			recent = append(recent, record)
			continue
		}
		
		endDate, err := time.Parse("2006-01-02T15:04:05", record.OfficeRecordEndDate)
		if err == nil && endDate.Year() >= 2023 {
			recent = append(recent, record)
		}
	}
	
	fmt.Printf("Records ending 2023 or later (or no end date): %d\n\n", len(recent))
	
	if len(recent) > 0 {
		fmt.Println("These recent records:")
		for _, record := range recent {
			endDate := record.OfficeRecordEndDate
			if endDate == "" {
				endDate = "CURRENT"
			}
			fmt.Printf("  - %s: %s → %s\n", record.OfficeRecordFullName, record.OfficeRecordStartDate, endDate)
		}
	}
}
