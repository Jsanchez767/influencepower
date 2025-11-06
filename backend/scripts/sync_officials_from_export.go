package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	log.Println("🏛️  Syncing Chicago officials from City API export endpoint...")
	log.Println("   Using /export/person with filter for 'Full City Council'")
	log.Println()

	// Get Chicago jurisdiction ID
	var jurisdictions []map[string]interface{}
	_, err := supabase.From("jurisdictions").
		Select("id", "exact", false).
		Eq("name", "Chicago").
		Eq("jurisdiction_type", "city").
		ExecuteTo(&jurisdictions)

	if err != nil || len(jurisdictions) == 0 {
		log.Fatal("Failed to find Chicago jurisdiction:", err)
	}

	chicagoID := int(jurisdictions[0]["id"].(float64))
	log.Printf("📍 Chicago jurisdiction ID: %d\n", chicagoID)

	// Create City API client
	client := cityapi.NewClient()

	// Get all Full City Council members using export endpoint
	log.Println("📋 Fetching current City Council members from export API...")
	
	// Use filter to get only Full City Council members
	persons, err := client.GetExportPersons("PersonType eq 'Full City Council'", "")
	if err != nil {
		log.Fatal("Failed to fetch persons:", err)
	}

	log.Printf("👤 Found %d Full City Council members\n\n", len(persons))

	if len(persons) == 0 {
		log.Println("⚠️  No current officials found in API")
		log.Println("💡 Try searching for specific names:")
		log.Println("   Example: filter=ward eq '14'")
		return
	}

	synced := 0
	created := 0
	updated := 0

	for _, person := range persons {
		// Extract ward from email or other fields
		ward := extractWard(person)
		
		if ward == 0 {
			log.Printf("⚠️  Could not determine ward for %s (skipping)", person.PersonFullName)
			continue
		}

		log.Printf("Processing Ward %d: %s...", ward, person.PersonFullName)

		// 1. Get or create person
		var existingPeople []map[string]interface{}
		_, err := supabase.From("people").
			Select("id", "exact", false).
			Eq("full_name", person.PersonFullName).
			ExecuteTo(&existingPeople)

		var personID int

		externalIDs := map[string]interface{}{
			"legistar_id":   person.PersonID,
			"legistar_guid": person.PersonGUID,
		}
		externalIDsJSON, _ := json.Marshal(externalIDs)

		if err != nil || len(existingPeople) == 0 {
			// Create new person
			personData := map[string]interface{}{
				"first_name":  person.PersonFirstName,
				"last_name":   person.PersonLastName,
				"full_name":   person.PersonFullName,
				"email":       person.PersonEmail,
				"phone":       person.PersonPhone,
				"website":     person.PersonWWW,
				"external_ids": string(externalIDsJSON),
			}

			var newPerson []map[string]interface{}
			_, err = supabase.From("people").
				Insert(personData, false, "", "", "").
				ExecuteTo(&newPerson)

			if err != nil {
				log.Printf("  ❌ Error creating person: %v", err)
				continue
			}

			personID = int(newPerson[0]["id"].(float64))
			log.Printf("  ✅ Created person ID: %d", personID)
			created++
		} else {
			personID = int(existingPeople[0]["id"].(float64))
			
			// Update person data
			updateData := map[string]interface{}{
				"email":       person.PersonEmail,
				"phone":       person.PersonPhone,
				"website":     person.PersonWWW,
				"external_ids": string(externalIDsJSON),
			}
			
			_, _, _ = supabase.From("people").
				Update(updateData, "", "").
				Eq("id", strconv.Itoa(personID)).
				Execute()
			
			log.Printf("  ℹ️  Updated person ID: %d", personID)
			updated++
		}

		// 2. Get position for this ward
		var positions []map[string]interface{}
		_, err = supabase.From("positions").
			Select("id", "exact", false).
			Eq("jurisdiction_id", strconv.Itoa(chicagoID)).
			Eq("position_type", "alderman").
			Eq("district_number", strconv.Itoa(ward)).
			ExecuteTo(&positions)

		if err != nil || len(positions) == 0 {
			log.Printf("  ❌ Position not found for ward %d", ward)
			continue
		}

		positionID := int(positions[0]["id"].(float64))

		// 3. Check if current term already exists
		var existingTerms []map[string]interface{}
		_, err = supabase.From("terms").
			Select("id", "exact", false).
			Eq("position_id", strconv.Itoa(positionID)).
			Eq("person_id", strconv.Itoa(personID)).
			ExecuteTo(&existingTerms)

		// Check if there's already a current term (end_date IS NULL)
		hasCurrentTerm := false
		for _, term := range existingTerms {
			if term["end_date"] == nil {
				hasCurrentTerm = true
				break
			}
		}

		if hasCurrentTerm {
			log.Printf("  ℹ️  Current term already exists")
			synced++
			continue
		}

		// 4. Create term (2023-2027 term started May 15, 2023)
		termData := map[string]interface{}{
			"position_id":   positionID,
			"person_id":     personID,
			"start_date":    "2023-05-15",
			"end_date":      nil, // Current term
			"term_number":   1,
			"election_type": "general",
		}

		var newTerm []map[string]interface{}
		_, err = supabase.From("terms").
			Insert(termData, false, "", "", "").
			ExecuteTo(&newTerm)

		if err != nil {
			log.Printf("  ❌ Error creating term: %v", err)
			continue
		}

		log.Printf("  ✅ Created current term")
		synced++
	}

	log.Println("\n🎉 Sync complete!")
	log.Printf("   Total synced: %d\n", synced)
	log.Printf("   People created: %d\n", created)
	log.Printf("   People updated: %d\n", updated)
	log.Println("\n💡 Next steps:")
	log.Println("   1. Verify: SELECT COUNT(*) FROM current_officials;")
	log.Println("   2. Run: go run scripts/update_headshots.go")
	log.Println("   3. Run: go run scripts/calculate_metrics.go")
}

// extractWard tries to extract ward number from various person fields
func extractWard(person cityapi.ExportPerson) int {
	// Try email first (Ward01@cityofchicago.org)
	if strings.Contains(person.PersonEmail, "Ward") {
		parts := strings.Split(person.PersonEmail, "@")
		if len(parts) > 0 {
			wardStr := strings.TrimPrefix(parts[0], "Ward")
			if ward, err := strconv.Atoi(wardStr); err == nil {
				return ward
			}
		}
	}

	// Try ward field if available
	if person.Ward > 0 {
		return person.Ward
	}

	return 0
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	log.Println(string(b))
}
