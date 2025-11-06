package main

import (
	"fmt"
	"log"
	"os"
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

	cityClient := cityapi.NewClient()

	log.Println("🖼️  Fetching headshots from City API...")

	// Get all persons from City API
	persons, err := cityClient.GetPersons()
	if err != nil {
		log.Fatalf("Failed to fetch persons: %v", err)
	}

	log.Printf("Found %d persons in City API", len(persons))

	// Get all officials from database
	var officials []map[string]interface{}
	_, err = supabase.From("officials").
		Select("id,name,city_api_person_id", "", false).
		ExecuteTo(&officials)

	if err != nil {
		log.Fatalf("Failed to fetch officials: %v", err)
	}

	log.Printf("Found %d officials in database", len(officials))

	updated := 0
	notFound := 0

	// Update each official's headshot
	for _, official := range officials {
		officialID := int(official["id"].(float64))
		officialName := official["name"].(string)
		
		var cityAPIPersonID string
		if official["city_api_person_id"] != nil {
			cityAPIPersonID = official["city_api_person_id"].(string)
		}

		log.Printf("\n📸 Processing %s (ID: %d)...", officialName, officialID)

		// Find matching person in City API
		var matchedPerson *cityapi.Person
		
		// First try by city_api_person_id if available
		if cityAPIPersonID != "" {
			for _, person := range persons {
				if person.PersonID == cityAPIPersonID {
					matchedPerson = &person
					break
				}
			}
		}

		// If not found, try by name matching
		if matchedPerson == nil {
			for _, person := range persons {
				if strings.EqualFold(person.PersonFullName, officialName) {
					matchedPerson = &person
					break
				}
			}
		}

		if matchedPerson == nil {
			log.Printf("  ⚠️  No match found in City API")
			notFound++
			continue
		}

		// Get the headshot URL from City API
		headshotURL := matchedPerson.PersonPhotoURL
		if headshotURL == "" {
			log.Printf("  ℹ️  No headshot URL available")
			continue
		}

		// Update the official in database
		updateData := map[string]interface{}{
			"image_url":            headshotURL,
			"city_api_person_id":   matchedPerson.PersonID,
		}

		var result []map[string]interface{}
		_, err := supabase.From("officials").
			Update(updateData, "", "").
			Eq("id", fmt.Sprintf("%d", officialID)).
			ExecuteTo(&result)

		if err != nil {
			log.Printf("  ❌ Failed to update: %v", err)
			continue
		}

		log.Printf("  ✅ Updated headshot URL: %s", headshotURL)
		updated++
	}

	log.Printf("\n\n✅ Headshot update complete!")
	log.Printf("   Updated: %d officials", updated)
	log.Printf("   Not found: %d officials", notFound)
}
