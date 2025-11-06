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
		
		var cityAPIPersonID int
		if official["city_api_person_id"] != nil {
			// Handle both string and numeric types
			switch v := official["city_api_person_id"].(type) {
			case string:
				fmt.Sscanf(v, "%d", &cityAPIPersonID)
			case float64:
				cityAPIPersonID = int(v)
			}
		}

		log.Printf("\n📸 Processing %s (ID: %d)...", officialName, officialID)

		// Find matching person in City API
		var matchedPerson *cityapi.Person
		
		// First try by city_api_person_id if available
		if cityAPIPersonID != 0 {
			for _, person := range persons {
				if person.PersonID == cityAPIPersonID {
					matchedPerson = &person
					break
				}
			}
		}

		// If not found, try by name matching
		if matchedPerson == nil {
			// Try exact match first
			for _, person := range persons {
				if strings.EqualFold(person.PersonFullName, officialName) {
					matchedPerson = &person
					break
				}
			}
			
			// If still not found, try matching with "LastName, FirstName" format
			if matchedPerson == nil {
				nameParts := strings.Fields(officialName)
				if len(nameParts) >= 2 {
					// Handle "FirstName LastName" -> "LastName, FirstName"
					lastName := nameParts[len(nameParts)-1]
					firstName := strings.Join(nameParts[:len(nameParts)-1], " ")
					reversedName := fmt.Sprintf("%s, %s", lastName, firstName)
					
					for _, person := range persons {
						if strings.EqualFold(person.PersonFullName, reversedName) {
							matchedPerson = &person
							break
						}
					}
				}
			}
		}

		if matchedPerson == nil {
			log.Printf("  ⚠️  No match found in City API")
			notFound++
			continue
		}

		// Construct the headshot URL from Legistar
		// The City API provides PersonGUID which can be used to construct image URLs
		headshotURL := fmt.Sprintf("https://chicago.legistar.com/people/%d", matchedPerson.PersonID)
		
		// Update the official in database
		updateData := map[string]interface{}{
			"image_url":            headshotURL,
			"city_api_person_id":   fmt.Sprintf("%d", matchedPerson.PersonID),
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
