package main

import (
	"fmt"
	"log"
	"os"
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

	log.Println("🔄 Calculating official metrics from City API data...")

	// Get all officials
	var officials []map[string]interface{}
	_, err := supabase.From("officials").
		Select("id,name,ward,party,role", "", false).
		ExecuteTo(&officials)

	if err != nil {
		log.Fatalf("Failed to fetch officials: %v", err)
	}

	log.Printf("Found %d officials to calculate metrics for", len(officials))

	// Calculate metrics for each official
	for _, official := range officials {
		officialID := int(official["id"].(float64))
		officialName := official["name"].(string)

		log.Printf("\n📊 Calculating metrics for %s (ID: %d)...", officialName, officialID)

		metrics := calculateOfficialMetrics(supabase, cityClient, officialID, officialName)
		
		if err := saveMetrics(supabase, officialID, metrics); err != nil {
			log.Printf("  ⚠️  Failed to save metrics: %v", err)
		} else {
			log.Printf("  ✓ Metrics saved successfully")
		}

		// Small delay to avoid overwhelming the API
		time.Sleep(200 * time.Millisecond)
	}

	log.Println("\n✅ Metrics calculation complete!")
}

type OfficialMetrics struct {
	// Legislative Productivity
	BillsIntroducedTotal       int
	BillsIntroducedCurrentTerm int
	BillsPassedTotal           int
	BillsPassedCurrentTerm     int
	FloorSpeechCount           int

	// Voting Record
	TotalVotesCast         int
	VotesYea               int
	VotesNay               int
	VotesPresent           int
	VotesAbsent            int
	VotingParticipationRate float64

	// Committee Attendance
	CommitteeAttendanceRate float64

	// Transparency Score
	TransparencyScore float64
}

func calculateOfficialMetrics(supabase *postgrest.Client, cityClient *cityapi.Client, officialID int, officialName string) OfficialMetrics {
	metrics := OfficialMetrics{}

	// 1. Calculate Legislative Productivity
	log.Printf("  📄 Analyzing legislation...")
	metrics.BillsIntroducedTotal, metrics.BillsIntroducedCurrentTerm = calculateLegislativeActivity(supabase, officialName)
	
	// 2. Calculate Voting Record
	log.Printf("  🗳️  Analyzing votes...")
	metrics.TotalVotesCast, metrics.VotesYea, metrics.VotesNay, metrics.VotesPresent = calculateVotingRecord(supabase, officialName)
	
	// Calculate voting participation (votes cast / total possible votes)
	// Estimate total possible votes from most active official
	if metrics.TotalVotesCast > 0 {
		// For now, assume 100% if they've voted
		metrics.VotingParticipationRate = 95.0 // Placeholder - would need total votes count
	}

	// 3. Calculate Committee Attendance
	log.Printf("  📋 Analyzing committee attendance...")
	metrics.CommitteeAttendanceRate = calculateCommitteeAttendance(supabase, officialName)

	// 4. Calculate Transparency Score (composite)
	log.Printf("  🔍 Calculating transparency score...")
	metrics.TransparencyScore = calculateTransparencyScore(metrics)

	return metrics
}

func calculateLegislativeActivity(supabase *postgrest.Client, officialName string) (int, int) {
	// Query matters where official is listed as sponsor
	var matters []map[string]interface{}
	_, err := supabase.From("matters").
		Select("matter_id,matter_sponsors,matter_intro_date,matter_status_name", "", false).
		ExecuteTo(&matters)

	if err != nil {
		log.Printf("    ⚠️  Error querying matters: %v", err)
		return 0, 0
	}

	totalIntroduced := 0
	currentTermIntroduced := 0
	currentTermStart := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC) // Brandon Johnson term start

	for _, matter := range matters {
		// Check if official is in sponsors (JSONB field)
		sponsorsStr, ok := matter["matter_sponsors"].(string)
		if !ok || sponsorsStr == "" {
			continue
		}

		// Simple name match in JSON string (proper JSON parsing would be better)
		if contains(sponsorsStr, officialName) {
			totalIntroduced++

			// Check if in current term
			if introDate, ok := matter["matter_intro_date"].(string); ok {
				if parsed, err := time.Parse(time.RFC3339, introDate); err == nil {
					if parsed.After(currentTermStart) {
						currentTermIntroduced++
					}
				}
			}
		}
	}

	log.Printf("    Bills introduced: %d total, %d current term", totalIntroduced, currentTermIntroduced)
	return totalIntroduced, currentTermIntroduced
}

func calculateVotingRecord(supabase *postgrest.Client, officialName string) (int, int, int, int) {
	// Query votes for this official
	var votes []map[string]interface{}
	_, err := supabase.From("votes").
		Select("vote_id,person_name,vote_value", "", false).
		Eq("person_name", officialName).
		ExecuteTo(&votes)

	if err != nil {
		log.Printf("    ⚠️  Error querying votes: %v", err)
		return 0, 0, 0, 0
	}

	totalVotes := len(votes)
	yea := 0
	nay := 0
	present := 0

	for _, vote := range votes {
		voteValue, ok := vote["vote_value"].(string)
		if !ok {
			continue
		}

		switch voteValue {
		case "Yea", "Yes", "Aye":
			yea++
		case "Nay", "No":
			nay++
		case "Present", "Abstain":
			present++
		}
	}

	log.Printf("    Votes: %d total (%d yea, %d nay, %d present)", totalVotes, yea, nay, present)
	return totalVotes, yea, nay, present
}

func calculateCommitteeAttendance(supabase *postgrest.Client, officialName string) float64 {
	// Query events (meetings) and check participation
	var events []map[string]interface{}
	_, err := supabase.From("events").
		Select("event_id,event_body_name,event_date", "", false).
		ExecuteTo(&events)

	if err != nil {
		log.Printf("    ⚠️  Error querying events: %v", err)
		return 0.0
	}

	// This is a simplified calculation
	// In a real system, we'd track actual attendance records
	// For now, we'll estimate based on vote participation
	if len(events) > 0 {
		// Placeholder: assume 85% attendance if they're active
		return 85.0
	}

	return 0.0
}

func calculateTransparencyScore(metrics OfficialMetrics) float64 {
	// Composite score based on multiple factors (0-100)
	score := 0.0

	// Voting participation (40% weight)
	score += metrics.VotingParticipationRate * 0.4

	// Committee attendance (30% weight)
	score += metrics.CommitteeAttendanceRate * 0.3

	// Legislative activity (20% weight)
	if metrics.BillsIntroducedCurrentTerm > 0 {
		// Award up to 20 points based on activity
		activityScore := float64(metrics.BillsIntroducedCurrentTerm) * 2.0
		if activityScore > 20.0 {
			activityScore = 20.0
		}
		score += activityScore
	}

	// Public engagement (10% weight) - placeholder
	score += 10.0 // Assume some baseline engagement

	log.Printf("    Transparency score: %.1f/100", score)
	return score
}

func saveMetrics(supabase *postgrest.Client, officialID int, metrics OfficialMetrics) error {
	metricsData := map[string]interface{}{
		"official_id":                   officialID,
		"bills_introduced_total":        metrics.BillsIntroducedTotal,
		"bills_introduced_current_term": metrics.BillsIntroducedCurrentTerm,
		"bills_passed_total":            metrics.BillsPassedTotal,
		"bills_passed_current_term":     metrics.BillsPassedCurrentTerm,
		"floor_speech_count":            metrics.FloorSpeechCount,
		"total_votes_cast":              metrics.TotalVotesCast,
		"votes_yea":                     metrics.VotesYea,
		"votes_nay":                     metrics.VotesNay,
		"votes_present":                 metrics.VotesPresent,
		"votes_absent":                  metrics.VotesAbsent,
		"voting_participation_rate":     metrics.VotingParticipationRate,
		"committee_attendance_rate":     metrics.CommitteeAttendanceRate,
		"transparency_score":            metrics.TransparencyScore,
		"last_calculated_at":            time.Now().Format(time.RFC3339),
	}

	_, _, err := supabase.From("official_metrics").Upsert(metricsData, "", "", "").Execute()
	return err
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
