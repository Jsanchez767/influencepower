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

	// Get all persons from City API
	persons, err := cityClient.GetPersons()
	if err != nil {
		log.Fatalf("Failed to fetch persons: %v", err)
	}

	// Get all officials from database
	var officials []map[string]interface{}
	_, err = supabase.From("officials").
		Select("id,name", "", false).
		ExecuteTo(&officials)

	if err != nil {
		log.Fatalf("Failed to fetch officials: %v", err)
	}

	fmt.Println("Unmatched officials and potential matches:\n")

	for _, official := range officials {
		officialName := official["name"].(string)
		
		// Try to find exact match
		found := false
		nameParts := strings.Fields(officialName)
		if len(nameParts) >= 2 {
			lastName := nameParts[len(nameParts)-1]
			firstName := strings.Join(nameParts[:len(nameParts)-1], " ")
			reversedName := fmt.Sprintf("%s, %s", lastName, firstName)
			
			for _, person := range persons {
				if strings.EqualFold(person.PersonFullName, reversedName) {
					found = true
					break
				}
			}
		}

		if !found {
			fmt.Printf("❌ %s\n", officialName)
			
			// Show similar names
			if len(nameParts) >= 2 {
				lastName := nameParts[len(nameParts)-1]
				
				fmt.Println("   Possible matches:")
				count := 0
				for _, person := range persons {
					if strings.Contains(strings.ToLower(person.PersonFullName), strings.ToLower(lastName)) {
						fmt.Printf("      - %s (ID: %d)\n", person.PersonFullName, person.PersonID)
						count++
						if count >= 3 {
							break
						}
					}
				}
			}
			fmt.Println()
		}
	}
}
