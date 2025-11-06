package main

import (
	"fmt"
	"log"
	"os"

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

	log.Println("🧹 Cleaning up database - removing all terms and people...")
	log.Println("⚠️  This will delete all synced data!")
	log.Println()

	// Delete all terms (will cascade to nothing since they're the leaf table)
	log.Println("1️⃣ Deleting all terms...")
	_, _, err := supabase.From("terms").
		Delete("", "").
		Neq("id", "0"). // Delete all rows
		Execute()
	
	if err != nil {
		log.Printf("   Error deleting terms: %v", err)
	} else {
		log.Println("   ✅ Terms deleted")
	}

	// Delete all people
	log.Println("2️⃣ Deleting all people...")
	_, _, err = supabase.From("people").
		Delete("", "").
		Neq("id", "0"). // Delete all rows
		Execute()
	
	if err != nil {
		log.Printf("   Error deleting people: %v", err)
	} else {
		log.Println("   ✅ People deleted")
	}

	// Delete all person_metrics
	log.Println("3️⃣ Deleting all person_metrics...")
	_, _, err = supabase.From("person_metrics").
		Delete("", "").
		Neq("person_id", "0"). // Delete all rows
		Execute()
	
	if err != nil {
		log.Printf("   Error deleting person_metrics: %v", err)
	} else {
		log.Println("   ✅ Person metrics deleted")
	}

	log.Println("\n🎉 Cleanup complete!")
	log.Println("\n💡 Next step: Run sync_from_elms_csv.go to repopulate with only current officials")
}
