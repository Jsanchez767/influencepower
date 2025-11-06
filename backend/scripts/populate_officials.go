package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
)

type Official struct {
	Name    string `json:"name"`
	Ward    *int   `json:"ward"`
	Party   string `json:"party"`
	Role    string `json:"role"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
}

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set")
	}

	client := postgrest.NewClient(supabaseURL+"/rest/v1", "", nil)
	client.TokenAuth(supabaseKey)

	// Read CSV file
	csvPath := "/Users/jesussanchez/Downloads/chicago_elected_officials_2025.csv"
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Skip header row
	for i, record := range records[1:] {
		position := record[0]
		name := record[1]
		wardStr := record[2]
		// office := record[3]
		party := record[4]
		// tookOffice := record[5]
		// communityAreas := record[6]
		contact := record[7]

		// Parse ward
		var ward *int
		if wardStr != "Citywide" && wardStr != "" {
			wardNum, err := strconv.Atoi(wardStr)
			if err == nil {
				ward = &wardNum
			}
		}

		// Determine role
		role := position
		if position == "Alderman" {
			role = "Alderman"
		} else if position == "Mayor" {
			role = "Mayor"
		} else if position == "City Clerk" {
			role = "City Clerk"
		} else if position == "City Treasurer" {
			role = "City Treasurer"
		}

		// Generate email if N/A
		email := contact
		if contact == "N/A" || contact == "" {
			if role == "Mayor" {
				email = "mayor@cityofchicago.org"
			} else if role == "City Clerk" {
				email = "cityclerk@cityofchicago.org"
			} else if role == "City Treasurer" {
				email = "treasurer@cityofchicago.org"
			} else if ward != nil {
				email = fmt.Sprintf("ward%02d@cityofchicago.org", *ward)
			}
		}

		official := Official{
			Name:    name,
			Ward:    ward,
			Party:   party,
			Role:    role,
			Contact: contact,
			Email:   email,
		}

		// Insert into Supabase
		var result []map[string]interface{}
		_, err = client.From("officials").Insert(official, false, "", "", "").ExecuteTo(&result)
		if err != nil {
			log.Printf("Failed to insert %s: %v", name, err)
		} else {
			log.Printf("✓ Inserted %s - %s", name, role)
		}

		if i > 0 && i%10 == 0 {
			log.Printf("Processed %d/%d officials...", i, len(records)-1)
		}
	}

	log.Printf("✓ Done! Inserted %d officials", len(records)-1)
}
