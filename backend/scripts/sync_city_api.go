package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jsanchez767/InfluencePower/backend/cityapi"
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
)

func main() {
	// Load environment variables
	godotenv.Load()

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Missing SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY environment variables")
	}

	// Initialize clients
	supabase := postgrest.NewClient(supabaseURL+"/rest/v1", "", map[string]string{
		"apikey":        supabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", supabaseKey),
	})

	cityClient := cityapi.NewClient()

	log.Println("🔄 Starting sync with Chicago City Clerk API...")

	// Sync bodies/committees first (needed for foreign keys)
	if err := syncBodies(supabase, cityClient); err != nil {
		log.Printf("⚠️  Error syncing bodies: %v", err)
	}

	// Sync persons (officials)
	if err := syncPersons(supabase, cityClient); err != nil {
		log.Printf("⚠️  Error syncing persons: %v", err)
	}

	// Sync recent matters (last 6 months)
	if err := syncRecentMatters(supabase, cityClient); err != nil {
		log.Printf("⚠️  Error syncing matters: %v", err)
	}

	// Sync recent events (last 3 months)
	if err := syncRecentEvents(supabase, cityClient); err != nil {
		log.Printf("⚠️  Error syncing events: %v", err)
	}

	log.Println("✅ Sync complete!")
}

func syncBodies(supabase *postgrest.Client, cityClient *cityapi.Client) error {
	log.Println("📋 Syncing bodies/committees...")

	bodies, err := cityClient.GetBodies()
	if err != nil {
		return fmt.Errorf("failed to fetch bodies: %w", err)
	}

	log.Printf("Found %d bodies", len(bodies))

	for _, body := range bodies {
		bodyData := map[string]interface{}{
			"body_id":        body.BodyID,
			"body_name":      body.BodyName,
			"body_type_id":   body.BodyTypeID,
			"body_type_name": body.BodyTypeName,
			"body_meet_flag": body.BodyMeetFlag,
		}

		_, _, err := supabase.From("bodies").Upsert(bodyData, "", "", "").Execute()
		if err != nil {
			log.Printf("  ⚠️  Failed to upsert body %s: %v", body.BodyName, err)
		} else {
			log.Printf("  ✓ Synced: %s", body.BodyName)
		}
	}

	return nil
}

func syncPersons(supabase *postgrest.Client, cityClient *cityapi.Client) error {
	log.Println("👥 Syncing persons...")

	persons, err := cityClient.GetPersons()
	if err != nil {
		return fmt.Errorf("failed to fetch persons: %w", err)
	}

	log.Printf("Found %d persons", len(persons))

	// Note: This is informational - we'll match these to our officials table later
	for _, person := range persons {
		log.Printf("  - %s (ID: %d)", person.PersonFullName, person.PersonID)
	}

	return nil
}

func syncRecentMatters(supabase *postgrest.Client, cityClient *cityapi.Client) error {
	log.Println("📄 Syncing recent matters...")

	// Try without date filter first to see what's available
	params := map[string]string{
		"$top": "100",
		"$orderby": "MatterIntroDate desc",
	}

	matters, err := cityClient.GetMatters(params)
	if err != nil {
		return fmt.Errorf("failed to fetch matters: %w", err)
	}

	log.Printf("Found %d recent matters", len(matters))

	for i, matter := range matters {
		if i >= 100 { // Limit to first 100 for initial sync
			log.Printf("  ⏸  Limiting to first 100 matters for initial sync")
			break
		}

		// Marshal sponsors and attachments to JSONB
		sponsorsJSON, _ := json.Marshal(matter.MatterSponsors)
		attachmentsJSON, _ := json.Marshal(matter.MatterAttachments)

		matterData := map[string]interface{}{
			"matter_id":              fmt.Sprintf("%d", matter.MatterID),
			"matter_file":            matter.MatterFile,
			"matter_name":            matter.MatterName,
			"matter_title":           matter.MatterTitle,
			"matter_type_id":         matter.MatterTypeID,
			"matter_type_name":       matter.MatterTypeName,
			"matter_status_id":       matter.MatterStatusID,
			"matter_status_name":     matter.MatterStatusName,
			"matter_intro_date":      parseAPIDate(matter.MatterIntroDate),
			"matter_agenda_date":     parseAPIDate(matter.MatterAgendaDate),
			"matter_passed_date":     parseAPIDate(matter.MatterPassedDate),
			"matter_enactment_date":  parseAPIDate(matter.MatterEnactmentDate),
			"matter_enactment_number": matter.MatterEnactmentNumber,
			"matter_requester":       matter.MatterRequester,
			"matter_sponsors":        string(sponsorsJSON),
			"matter_attachments":     string(attachmentsJSON),
			"matter_text":            matter.MatterText,
			"matter_version":         matter.MatterVersion,
		}

		_, _, err := supabase.From("matters").Upsert(matterData, "", "", "").Execute()
		if err != nil {
			log.Printf("  ⚠️  Failed to upsert matter %s: %v", matter.MatterFile, err)
		} else {
			log.Printf("  ✓ Synced: %s - %s", matter.MatterFile, truncate(matter.MatterTitle, 60))
		}

		// Fetch and sync votes for this matter
		syncVotesForMatter(supabase, cityClient, matter.MatterID)

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func syncVotesForMatter(supabase *postgrest.Client, cityClient *cityapi.Client, matterID int) {
	votes, err := cityClient.GetVotes(matterID)
	if err != nil {
		log.Printf("    ⚠️  Failed to fetch votes for matter %d: %v", matterID, err)
		return
	}

	if len(votes) == 0 {
		return
	}

	log.Printf("    📊 Found %d votes for matter %d", len(votes), matterID)

	for _, vote := range votes {
		voteData := map[string]interface{}{
			"vote_id":      fmt.Sprintf("%d", vote.VoteID),
			"matter_id":    fmt.Sprintf("%d", vote.VoteMatterID),
			"person_id":    vote.VotePersonID,
			"person_name":  vote.VotePersonName,
			"vote_value":   vote.VoteValue,
			"vote_date":    parseAPIDate(vote.VoteDate),
			"vote_event_id": vote.VoteEventID,
		}

		_, _, err := supabase.From("votes").Upsert(voteData, "", "", "").Execute()
		if err != nil {
			log.Printf("      ⚠️  Failed to upsert vote: %v", err)
		}
	}
}

func syncRecentEvents(supabase *postgrest.Client, cityClient *cityapi.Client) error {
	log.Println("📅 Syncing recent events...")

	// Try without date filter first
	params := map[string]string{
		"$top": "100",
		"$orderby": "EventDate desc",
	}

	events, err := cityClient.GetEvents(params)
	if err != nil {
		return fmt.Errorf("failed to fetch events: %w", err)
	}

	log.Printf("Found %d recent events", len(events))

	for i, event := range events {
		if i >= 50 { // Limit to first 50 for initial sync
			log.Printf("  ⏸  Limiting to first 50 events for initial sync")
			break
		}

		// Marshal event items to JSONB
		itemsJSON, _ := json.Marshal(event.EventItems)

		eventData := map[string]interface{}{
			"event_id":          fmt.Sprintf("%d", event.EventID),
			"event_body_id":     event.EventBodyID,
			"event_body_name":   event.EventBodyName,
			"event_date":        parseAPIDate(event.EventDate),
			"event_time":        event.EventTime,
			"event_location":    event.EventLocation,
			"event_agenda_file": event.EventAgendaFile,
			"event_minutes_file": event.EventMinutesFile,
			"event_video_url":   event.EventVideoURL,
			"event_items":       string(itemsJSON),
		}

		_, _, err := supabase.From("events").Upsert(eventData, "", "", "").Execute()
		if err != nil {
			log.Printf("  ⚠️  Failed to upsert event: %v", err)
		} else {
			log.Printf("  ✓ Synced: %s on %s", event.EventBodyName, event.EventDate[:10])
		}

		// Sync event items
		for _, item := range event.EventItems {
			itemData := map[string]interface{}{
				"event_item_id":          fmt.Sprintf("%d", item.EventItemID),
				"event_id":               fmt.Sprintf("%d", event.EventID),
				"matter_id":              fmt.Sprintf("%d", item.EventItemMatterID),
				"item_agenda_sequence":   item.EventItemAgendaSequence,
				"item_agenda_number":     item.EventItemAgendaNumber,
				"item_action":            item.EventItemAction,
				"item_action_text":       item.EventItemActionText,
			}

			_, _, err := supabase.From("event_items").Upsert(itemData, "", "", "").Execute()
			if err != nil {
				log.Printf("    ⚠️  Failed to upsert event item: %v", err)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// parseAPIDate parses the API date format and returns a timestamp
func parseAPIDate(dateStr string) *string {
	if dateStr == "" {
		return nil
	}

	// Chicago API uses format like "2024-01-15T00:00:00"
	parsed, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try alternate format
		parsed, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return nil
		}
	}

	formatted := parsed.Format(time.RFC3339)
	return &formatted
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
