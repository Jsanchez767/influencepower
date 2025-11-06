package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jsanchez767/InfluencePower/backend/cityapi"
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
)

func main() {
	godotenv.Load()

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Missing SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY environment variables")
	}

	supabase := postgrest.NewClient(supabaseURL+"/rest/v1", "", map[string]string{
		"apikey":        supabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", supabaseKey),
	})

	cityClient := cityapi.NewClient()

	log.Println("🏛️  Syncing Chicago officials from City API...")
	log.Println("   Using position-based matching (jurisdiction + district)")
	log.Println()

	// Get Chicago jurisdiction ID
	var jurisdictions []map[string]interface{}
	_, err := supabase.From("jurisdictions").
		Select("id", "", false).
		Eq("name", "Chicago").
		Eq("jurisdiction_type", "city").
		ExecuteTo(&jurisdictions)

	if err != nil || len(jurisdictions) == 0 {
		log.Fatal("Chicago jurisdiction not found. Please run migration first.")
	}

	chicagoID := int(jurisdictions[0]["id"].(float64))
	log.Printf("📍 Chicago jurisdiction ID: %d\n\n", chicagoID)

	// Fetch all office records from City API
	officeRecords, err := cityClient.GetOfficeRecords()
	if err != nil {
		log.Fatalf("Failed to fetch office records: %v", err)
	}

	log.Printf("📋 Found %d office records\n", len(officeRecords))

	// Fetch all persons for headshots
	persons, err := cityClient.GetPersons()
	if err != nil {
		log.Fatalf("Failed to fetch persons: %v", err)
	}

	log.Printf("👤 Found %d persons\n\n", len(persons))

	// Build person lookup map
	personMap := make(map[int]cityapi.Person)
	for _, p := range persons {
		personMap[p.PersonID] = p
	}

	// Filter for current aldermen and citywide officials
	currentOfficials := filterCurrentOfficials(officeRecords)
	log.Printf("✅ Found %d current officials\n\n", len(currentOfficials))

	synced := 0
	updated := 0
	created := 0

	for _, record := range currentOfficials {
		// Extract ward from email or title
		ward := extractWard(record)
		positionType := determinePositionType(record, ward)

		log.Printf("🔄 Processing: %s - %s (Ward %v)\n",
			record.OfficeRecordFullName,
			record.OfficeRecordTitle,
			ward)

		// Get person details for headshot
		person, hasPerson := personMap[record.OfficeRecordPersonID]
		imageURL := ""
		if hasPerson {
			imageURL = fmt.Sprintf("https://chicago.legistar.com/people/%d", person.PersonID)
		}

		// 1. Find position by jurisdiction + district
		var positions []map[string]interface{}
		query := supabase.From("positions").
			Select("id", "", false).
			Eq("jurisdiction_id", fmt.Sprintf("%d", chicagoID)).
			Eq("position_type", positionType)

		if ward > 0 {
			query = query.Eq("district_number", fmt.Sprintf("%d", ward))
		} else {
			query = query.Is("district_number", "null")
		}

		_, err := query.ExecuteTo(&positions)

		if err != nil {
			log.Printf("   ❌ Failed to query positions: %v\n", err)
			continue
		}

		var positionID int
		if len(positions) == 0 {
			log.Printf("   ⚠️  Position not found, skipping\n")
			continue
		} else {
			positionID = int(positions[0]["id"].(float64))
		}

		// 2. Upsert person with external_ids JSONB
		var people []map[string]interface{}
		_, err = supabase.From("people").
			Select("id", "", false).
			Eq("external_ids->>legistar_id", fmt.Sprintf("%d", record.OfficeRecordPersonID)).
			ExecuteTo(&people)

		if err != nil {
			log.Printf("   ❌ Failed to query people: %v\n", err)
			continue
		}

		var personID int
		if len(people) == 0 {
			// Create new person
			externalIDs := map[string]interface{}{
				"legistar_id":   record.OfficeRecordPersonID,
				"legistar_guid": "",
			}

			if hasPerson {
				externalIDs["legistar_guid"] = person.PersonGUID
			}

			personData := map[string]interface{}{
				"external_ids": externalIDs,
				"first_name":   record.OfficeRecordFirstName,
				"last_name":    record.OfficeRecordLastName,
				"full_name":    record.OfficeRecordFullName,
				"email":        record.OfficeRecordEmail,
				"image_url":    imageURL,
			}

			if hasPerson {
				personData["phone"] = person.PersonPhone
				personData["website"] = person.PersonWWW
				personData["headshot_last_updated"] = time.Now().Format(time.RFC3339)
			}

			var result []map[string]interface{}
			_, err = supabase.From("people").
				Insert(personData, false, "", "", "").
				ExecuteTo(&result)

			if err != nil {
				log.Printf("   ❌ Failed to create person: %v\n", err)
				continue
			}

			personID = int(result[0]["id"].(float64))
			log.Printf("   ✨ Created person: %s\n", record.OfficeRecordFullName)
			created++
		} else {
			personID = int(people[0]["id"].(float64))

			// Update person data with external_ids
			externalIDs := map[string]interface{}{
				"legistar_id": record.OfficeRecordPersonID,
			}

			if hasPerson {
				externalIDs["legistar_guid"] = person.PersonGUID
			}

			updateData := map[string]interface{}{
				"external_ids": externalIDs,
				"full_name":    record.OfficeRecordFullName,
				"first_name":   record.OfficeRecordFirstName,
				"last_name":    record.OfficeRecordLastName,
				"email":        record.OfficeRecordEmail,
				"image_url":    imageURL,
			}

			if hasPerson {
				updateData["phone"] = person.PersonPhone
				updateData["website"] = person.PersonWWW
				updateData["headshot_last_updated"] = time.Now().Format(time.RFC3339)
			}

			var result []map[string]interface{}
			_, err = supabase.From("people").
				Update(updateData, "", "").
				Eq("id", fmt.Sprintf("%d", personID)).
				ExecuteTo(&result)

			if err != nil {
				log.Printf("   ❌ Failed to update person: %v\n", err)
			} else {
				log.Printf("   📝 Updated person\n")
				updated++
			}
		}

		// 3. Create or update term
		startDate := parseDate(record.OfficeRecordStartDate)
		endDate := parseDate(record.OfficeRecordEndDate)

		var terms []map[string]interface{}
		_, err = supabase.From("terms").
			Select("id", "", false).
			Eq("external_id", fmt.Sprintf("%d", record.OfficeRecordID)).
			ExecuteTo(&terms)

		if err != nil {
			log.Printf("   ❌ Failed to query terms: %v\n", err)
			continue
		}

		if len(terms) == 0 {
			// Create new term
			termData := map[string]interface{}{
				"position_id":   positionID,
				"person_id":     personID,
				"start_date":    startDate,
				"external_id":   record.OfficeRecordID,
				"external_guid": record.OfficeRecordGUID,
			}

			if endDate != "" {
				termData["end_date"] = endDate
			}

			var result []map[string]interface{}
			_, err = supabase.From("terms").
				Insert(termData, false, "", "", "").
				ExecuteTo(&result)

			if err != nil {
				log.Printf("   ❌ Failed to create term: %v\n", err)
			} else {
				log.Printf("   ✅ Created term\n")
			}
		} else {
			// Update term
			termData := map[string]interface{}{
				"person_id":  personID,
				"start_date": startDate,
			}

			if endDate != "" {
				termData["end_date"] = endDate
			}

			var result []map[string]interface{}
			_, err = supabase.From("terms").
				Update(termData, "", "").
				Eq("id", fmt.Sprintf("%d", int(terms[0]["id"].(float64)))).
				ExecuteTo(&result)

			if err != nil {
				log.Printf("   ❌ Failed to update term: %v\n", err)
			} else {
				log.Printf("   ✅ Updated term\n")
			}
		}

		synced++
		log.Println()
	}

	log.Printf("\n🎉 Sync complete!\n")
	log.Printf("   Total synced: %d\n", synced)
	log.Printf("   Created: %d\n", created)
	log.Printf("   Updated: %d\n", updated)
}

func filterCurrentOfficials(records []cityapi.OfficeRecord) []cityapi.OfficeRecord {
	var current []cityapi.OfficeRecord
	now := time.Now()

	for _, record := range records {
		// Check if this is a city council position
		if record.OfficeRecordBodyName != "City Council" {
			continue
		}

		// Check if currently serving (no end date or end date in future)
		if record.OfficeRecordEndDate == "" {
			current = append(current, record)
			continue
		}

		endDate, err := time.Parse("2006-01-02T15:04:05", record.OfficeRecordEndDate)
		if err == nil && endDate.After(now) {
			current = append(current, record)
			continue
		}
	}

	return current
}

func extractWard(record cityapi.OfficeRecord) int {
	// Try to extract ward number from email or extra text
	// Email format: Ward01@cityofchicago.org
	re := regexp.MustCompile(`[Ww]ard\s*0*(\d+)`)

	if matches := re.FindStringSubmatch(record.OfficeRecordEmail); len(matches) > 1 {
		if ward, err := strconv.Atoi(matches[1]); err == nil {
			return ward
		}
	}

	if matches := re.FindStringSubmatch(record.OfficeRecordExtraText); len(matches) > 1 {
		if ward, err := strconv.Atoi(matches[1]); err == nil {
			return ward
		}
	}

	return 0
}

func determinePositionType(record cityapi.OfficeRecord, ward int) string {
	title := strings.ToLower(record.OfficeRecordTitle)

	if strings.Contains(title, "mayor") {
		return "mayor"
	}
	if strings.Contains(title, "clerk") {
		return "clerk"
	}
	if strings.Contains(title, "treasurer") {
		return "treasurer"
	}
	if ward > 0 {
		return "alderman"
	}

	return "other"
}

func parseDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Parse ISO format: 2023-05-14T00:00:00
	t, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		return ""
	}

	// Return as YYYY-MM-DD
	return t.Format("2006-01-02")
}
