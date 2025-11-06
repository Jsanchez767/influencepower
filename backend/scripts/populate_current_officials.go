package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Jsanchez767/InfluencePower/backend/db"
)

// Current Chicago officials (2023-2027 term)
var currentOfficials = []struct {
	Ward      int
	FirstName string
	LastName  string
	Email     string
	Party     string
}{
	{1, "Daniel", "La Spata", "Ward01@cityofchicago.org", "Democratic"},
	{2, "Brian", "Hopkins", "Ward02@cityofchicago.org", "Democratic"},
	{3, "Pat", "Dowell", "Ward03@cityofchicago.org", "Democratic"},
	{4, "Lamont", "Robinson", "Ward04@cityofchicago.org", "Democratic"},
	{5, "Desmon", "Yancy", "Ward05@cityofchicago.org", "Democratic"},
	{6, "William", "Hall", "Ward06@cityofchicago.org", "Democratic"},
	{7, "Greg", "Mitchell", "Ward07@cityofchicago.org", "Democratic"},
	{8, "Michelle", "Harris", "Ward08@cityofchicago.org", "Democratic"},
	{9, "Anthony", "Beale", "Ward09@cityofchicago.org", "Democratic"},
	{10, "Peter", "Chico", "Ward10@cityofchicago.org", "Democratic"},
	{11, "Nicole", "Lee", "Ward11@cityofchicago.org", "Democratic"},
	{12, "Julia", "Ramirez", "Ward12@cityofchicago.org", "Democratic"},
	{13, "Marty", "Quinn", "Ward13@cityofchicago.org", "Democratic"},
	{14, "Jeannette", "Taylor", "Ward14@cityofchicago.org", "Democratic"},
	{15, "Raymond", "Lopez", "Ward15@cityofchicago.org", "Democratic"},
	{16, "Stephanie", "Coleman", "Ward16@cityofchicago.org", "Democratic"},
	{17, "David", "Moore", "Ward17@cityofchicago.org", "Democratic"},
	{18, "Derrick", "Curtis", "Ward18@cityofchicago.org", "Democratic"},
	{19, "Matthew", "O'Shea", "Ward19@cityofchicago.org", "Democratic"},
	{20, "Jeanette", "Taylor", "Ward20@cityofchicago.org", "Democratic"},
	{21, "Ronnie", "Mosley", "Ward21@cityofchicago.org", "Democratic"},
	{22, "Michael", "Rodriguez", "Ward22@cityofchicago.org", "Democratic"},
	{23, "Silvana", "Tabares", "Ward23@cityofchicago.org", "Democratic"},
	{24, "Monique", "Scott", "Ward24@cityofchicago.org", "Democratic"},
	{25, "Byron", "Sigcho-Lopez", "Ward25@cityofchicago.org", "Democratic"},
	{26, "Jessie", "Fuentes", "Ward26@cityofchicago.org", "Democratic"},
	{27, "Walter", "Burnett", "Ward27@cityofchicago.org", "Democratic"},
	{28, "Jason", "Ervin", "Ward28@cityofchicago.org", "Democratic"},
	{29, "Chris", "Taliaferro", "Ward29@cityofchicago.org", "Democratic"},
	{30, "Ruth", "Cruz", "Ward30@cityofchicago.org", "Democratic"},
	{31, "Felix", "Cardona", "Ward31@cityofchicago.org", "Democratic"},
	{32, "Scott", "Waguespack", "Ward32@cityofchicago.org", "Democratic"},
	{33, "Rossana", "Rodriguez-Sanchez", "Ward33@cityofchicago.org", "Democratic"},
	{34, "Bill", "Conway", "Ward34@cityofchicago.org", "Democratic"},
	{35, "Carlos", "Ramirez-Rosa", "Ward35@cityofchicago.org", "Democratic"},
	{36, "Gilbert", "Villegas", "Ward36@cityofchicago.org", "Democratic"},
	{37, "Emma", "Mitts", "Ward37@cityofchicago.org", "Democratic"},
	{38, "Nick", "Sposato", "Ward38@cityofchicago.org", "Democratic"},
	{39, "Sam", "Nugent", "Ward39@cityofchicago.org", "Democratic"},
	{40, "Andre", "Vasquez", "Ward40@cityofchicago.org", "Democratic"},
	{41, "Anthony", "Napolitano", "Ward41@cityofchicago.org", "Republican"},
	{42, "Brendan", "Reilly", "Ward42@cityofchicago.org", "Democratic"},
	{43, "Timmy", "Knudsen", "Ward43@cityofchicago.org", "Democratic"},
	{44, "Bennett", "Lawson", "Ward44@cityofchicago.org", "Democratic"},
	{45, "Jim", "Gardiner", "Ward45@cityofchicago.org", "Republican"},
	{46, "Angela", "Clay", "Ward46@cityofchicago.org", "Democratic"},
	{47, "Matt", "Martin", "Ward47@cityofchicago.org", "Democratic"},
	{48, "Leni", "Manaa-Hoppenworth", "Ward48@cityofchicago.org", "Democratic"},
	{49, "Maria", "Hadden", "Ward49@cityofchicago.org", "Democratic"},
	{50, "Debra", "Silverstein", "Ward50@cityofchicago.org", "Democratic"},
}

func main() {
	log.Println("🏛️  Populating current Chicago officials (2023-2027 term)...")

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Get Chicago jurisdiction ID
	var jurisdictions []map[string]interface{}
	_, err := db.Client.From("jurisdictions").
		Select("id", "exact", false).
		Eq("name", "Chicago").
		Eq("jurisdiction_type", "city").
		ExecuteTo(&jurisdictions)

	if err != nil || len(jurisdictions) == 0 {
		log.Fatal("Failed to find Chicago jurisdiction:", err)
	}

	chicagoID := int(jurisdictions[0]["id"].(float64))
	log.Printf("📍 Chicago jurisdiction ID: %d\n", chicagoID)

	created := 0
	updated := 0

	for _, official := range currentOfficials {
		fullName := official.FirstName + " " + official.LastName

		log.Printf("Processing Ward %d: %s...", official.Ward, fullName)

		// 1. Get or create person
		var existingPeople []map[string]interface{}
		_, err := db.Client.From("people").
			Select("id", "exact", false).
			Eq("full_name", fullName).
			ExecuteTo(&existingPeople)

		var personID int

		if err != nil || len(existingPeople) == 0 {
			// Create new person
			personData := map[string]interface{}{
				"first_name":       official.FirstName,
				"last_name":        official.LastName,
				"full_name":        fullName,
				"email":            official.Email,
				"party_affiliation": official.Party,
			}

			var newPerson []map[string]interface{}
			_, err = db.Client.From("people").
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
			log.Printf("  ℹ️  Person already exists (ID: %d)", personID)
			
			// Update email and party if needed
			updateData := map[string]interface{}{
				"email":            official.Email,
				"party_affiliation": official.Party,
			}
			
			_, _ = db.Client.From("people").
				Update(updateData, "", "").
				Eq("id", personID).
				Execute()
			
			updated++
		}

		// 2. Get position for this ward
		var positions []map[string]interface{}
		_, err = db.Client.From("positions").
			Select("id", "exact", false).
			Eq("jurisdiction_id", chicagoID).
			Eq("position_type", "alderman").
			Eq("district_number", official.Ward).
			ExecuteTo(&positions)

		if err != nil || len(positions) == 0 {
			log.Printf("  ❌ Position not found for ward %d", official.Ward)
			continue
		}

		positionID := int(positions[0]["id"].(float64))

		// 3. Check if term already exists
		var existingTerms []map[string]interface{}
		_, err = db.Client.From("terms").
			Select("id", "exact", false).
			Eq("position_id", positionID).
			Eq("person_id", personID).
			ExecuteTo(&existingTerms)

		if err == nil && len(existingTerms) > 0 {
			log.Printf("  ℹ️  Term already exists")
			continue
		}

		// 4. Create term
		termData := map[string]interface{}{
			"position_id": positionID,
			"person_id":   personID,
			"start_date":  "2023-05-15",
			"end_date":    nil, // Current term
			"term_number": 1,
			"election_type": "general",
		}

		var newTerm []map[string]interface{}
		_, err = db.Client.From("terms").
			Insert(termData, false, "", "", "").
			ExecuteTo(&newTerm)

		if err != nil {
			log.Printf("  ❌ Error creating term: %v", err)
			continue
		}

		log.Printf("  ✅ Created term")
	}

	log.Println("\n🎉 Population complete!")
	log.Printf("   People created: %d\n", created)
	log.Printf("   People updated: %d\n", updated)
	log.Println("\n💡 Next steps:")
	log.Println("   1. Verify: SELECT COUNT(*) FROM current_officials;")
	log.Println("   2. Run: go run scripts/update_headshots.go")
	log.Println("   3. Run: go run scripts/calculate_metrics.go")
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	log.Println(string(b))
}
