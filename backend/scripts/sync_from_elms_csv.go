package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	log.Println("🏛️  Syncing Chicago officials from City Clerk ELMS API...")
	log.Println("   Using CSV export endpoint")
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

	// Fetch CSV data from ELMS API
	log.Println("📋 Fetching current officials from ELMS API...")
	url := "https://api.chicityclerkelms.chicago.gov/export/person"
	
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to fetch data:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("API returned status %d", resp.StatusCode)
	}

	// Parse CSV
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to parse CSV:", err)
	}

	log.Printf("👤 Found %d officials (including header)\n\n", len(records))

	if len(records) < 2 {
		log.Fatal("No data in CSV")
	}

	// First row is header
	// Name,Office,Office Phone,Fax,Email,Website,Office Address,City,State,Zip,City Hall Phone,City Hall Address,City,State,Zip
	
	synced := 0
	created := 0
	updated := 0
	skipped := 0

	for i, record := range records[1:] { // Skip header
		if len(record) < 5 {
			log.Printf("⚠️  Row %d: Insufficient columns, skipping", i+2)
			skipped++
			continue
		}

		name := strings.TrimSpace(record[0])
		wardStr := strings.TrimSpace(record[1])
		phone := strings.TrimSpace(record[2])
		email := strings.TrimSpace(record[4])
		website := strings.TrimSpace(record[5])

		// Parse name (format: "Last, First Middle")
		nameParts := strings.Split(name, ",")
		if len(nameParts) < 2 {
			log.Printf("⚠️  Could not parse name: %s, skipping", name)
			skipped++
			continue
		}

		lastName := strings.TrimSpace(strings.Trim(nameParts[0], "\""))
		firstMiddle := strings.TrimSpace(strings.Trim(nameParts[1], "\""))
		firstNameParts := strings.Fields(firstMiddle)
		firstName := firstNameParts[0]
		fullName := firstName + " " + lastName

		// Parse ward number
		ward, err := strconv.Atoi(wardStr)
		if err != nil {
			log.Printf("⚠️  Invalid ward '%s' for %s, skipping", wardStr, fullName)
			skipped++
			continue
		}

		if ward < 1 || ward > 50 {
			log.Printf("⚠️  Ward %d out of range for %s, skipping", ward, fullName)
			skipped++
			continue
		}

		log.Printf("Processing Ward %d: %s...", ward, fullName)

		// 1. Get or create person
		var existingPeople []map[string]interface{}
		_, err = supabase.From("people").
			Select("id", "exact", false).
			Eq("full_name", fullName).
			ExecuteTo(&existingPeople)

		var personID int

		if err != nil || len(existingPeople) == 0 {
			// Create new person
			personData := map[string]interface{}{
				"first_name": firstName,
				"last_name":  lastName,
				"full_name":  fullName,
				"email":      email,
				"phone":      phone,
				"website":    website,
			}

			var newPerson []map[string]interface{}
			_, err = supabase.From("people").
				Insert(personData, false, "", "", "").
				ExecuteTo(&newPerson)

			if err != nil {
				log.Printf("  ❌ Error creating person: %v", err)
				skipped++
				continue
			}

			personID = int(newPerson[0]["id"].(float64))
			log.Printf("  ✅ Created person ID: %d", personID)
			created++
		} else {
			personID = int(existingPeople[0]["id"].(float64))
			
			// Update contact info
			updateData := map[string]interface{}{
				"email":   email,
				"phone":   phone,
				"website": website,
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
			skipped++
			continue
		}

		positionID := int(positions[0]["id"].(float64))

		// 3. Check if current term already exists
		var existingTerms []map[string]interface{}
		_, err = supabase.From("terms").
			Select("id, end_date", "exact", false).
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

		// 4. Create current term (2023-2027 election started May 15, 2023)
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
			skipped++
			continue
		}

		log.Printf("  ✅ Created current term")
		synced++
	}

	log.Println("\n🎉 Sync complete!")
	log.Printf("   Total synced: %d\n", synced)
	log.Printf("   People created: %d\n", created)
	log.Printf("   People updated: %d\n", updated)
	log.Printf("   Skipped: %d\n", skipped)
	log.Println("\n💡 Next steps:")
	log.Println("   1. Verify: SELECT COUNT(*) FROM current_officials;")
	log.Println("   2. Run: go run scripts/update_headshots.go")
	log.Println("   3. Run: go run scripts/calculate_metrics.go")
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	log.Println(string(b))
}

// extractWard extracts ward number from email or other fields
func extractWard(email string) int {
	// Match Ward01, Ward02, etc.
	re := regexp.MustCompile(`Ward(\d+)`)
	matches := re.FindStringSubmatch(email)
	if len(matches) > 1 {
		if ward, err := strconv.Atoi(matches[1]); err == nil {
			return ward
		}
	}
	return 0
}
