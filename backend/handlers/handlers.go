package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jsanchez767/InfluencePower/backend/db"
	"github.com/Jsanchez767/InfluencePower/backend/models"
	"github.com/gorilla/mux"
)

// GetOfficials returns all officials
func GetOfficials(w http.ResponseWriter, r *http.Request) {
	var officials []models.Official
	
	// Query all officials from Supabase
	_, err := db.Client.From("officials").
		Select("*", "exact", false).
		ExecuteTo(&officials)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(officials)
}

// GetOfficialByID returns a single official by ID
func GetOfficialByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var officials []models.Official
	
	_, err := db.Client.From("officials").
		Select("*", "exact", false).
		Eq("id", id).
		ExecuteTo(&officials)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(officials) == 0 {
		http.Error(w, "Official not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(officials[0])
}

// CreateOfficial creates a new official
func CreateOfficial(w http.ResponseWriter, r *http.Request) {
	var official models.Official
	
	if err := json.NewDecoder(r.Body).Decode(&official); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result []models.Official
	_, err := db.Client.From("officials").
		Insert(official, false, "", "", "").
		ExecuteTo(&result)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if len(result) > 0 {
		json.NewEncoder(w).Encode(result[0])
	} else {
		json.NewEncoder(w).Encode(official)
	}
}

// UpdateOfficial updates an existing official
func UpdateOfficial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var official models.Official
	if err := json.NewDecoder(r.Body).Decode(&official); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result []models.Official
	_, err := db.Client.From("officials").
		Update(official, "", "").
		Eq("id", id).
		ExecuteTo(&result)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(result) > 0 {
		json.NewEncoder(w).Encode(result[0])
	} else {
		json.NewEncoder(w).Encode(official)
	}
}

// DeleteOfficial deletes an official
func DeleteOfficial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var result []models.Official
	_, err := db.Client.From("officials").
		Delete("", "").
		Eq("id", id).
		ExecuteTo(&result)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetOfficialsByParty returns officials filtered by party
func GetOfficialsByParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	party := vars["party"]

	var officials []models.Official
	
	_, err := db.Client.From("officials").
		Select("*", "exact", false).
		Eq("party", party).
		ExecuteTo(&officials)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(officials)
}

// GetOfficialsByWard returns the official for a specific ward
func GetOfficialsByWard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ward := vars["ward"]

	var officials []models.Official
	
	_, err := db.Client.From("officials").
		Select("*", "exact", false).
		Eq("ward", ward).
		ExecuteTo(&officials)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(officials)
}

// GetVotingRecords returns voting records for an official
func GetVotingRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	officialID := vars["id"]

	var records []models.VotingRecord
	
	_, err := db.Client.From("voting_records").
		Select("*", "exact", false).
		Eq("official_id", officialID).
		Order("vote_date", nil).
		ExecuteTo(&records)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// CreateVotingRecord creates a new voting record
func CreateVotingRecord(w http.ResponseWriter, r *http.Request) {
	var record models.VotingRecord
	
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result []models.VotingRecord
	_, err := db.Client.From("voting_records").
		Insert(record, false, "", "", "").
		ExecuteTo(&result)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if len(result) > 0 {
		json.NewEncoder(w).Encode(result[0])
	} else {
		json.NewEncoder(w).Encode(record)
	}
}

// GetWardStatistics returns statistics for a specific ward
func GetWardStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wardStr := vars["ward"]
	ward, err := strconv.Atoi(wardStr)
	if err != nil {
		http.Error(w, "Invalid ward number", http.StatusBadRequest)
		return
	}

	var stats []models.WardStatistic
	
	_, err = db.Client.From("ward_statistics").
		Select("*", "exact", false).
		Eq("ward", strconv.Itoa(ward)).
		ExecuteTo(&stats)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(stats) == 0 {
		http.Error(w, "Ward statistics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats[0])
}

// GetCommittees returns all committees
func GetCommittees(w http.ResponseWriter, r *http.Request) {
	var committees []models.Committee
	
	_, err := db.Client.From("committees").
		Select("*", "exact", false).
		ExecuteTo(&committees)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(committees)
}

// GetOfficialCommittees returns committees for a specific official
func GetOfficialCommittees(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	officialID := vars["id"]

	var committees []models.OfficialCommittee
	
	_, err := db.Client.From("official_committees").
		Select("*, committees(*)", "exact", false).
		Eq("official_id", officialID).
		ExecuteTo(&committees)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(committees)
}

// GetOfficialMetrics returns performance metrics for a specific official
func GetOfficialMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	officialID := vars["id"]

	var metrics []map[string]interface{}
	
	_, err := db.Client.From("official_metrics").
		Select("*", "exact", false).
		Eq("official_id", officialID).
		ExecuteTo(&metrics)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(metrics) == 0 {
		http.Error(w, "Metrics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics[0])
}

// GetWardMetrics returns metrics for a specific ward
func GetWardMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ward := vars["ward"]

	var metrics []map[string]interface{}
	
	_, err := db.Client.From("ward_metrics").
		Select("*", "exact", false).
		Eq("ward", ward).
		ExecuteTo(&metrics)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(metrics) == 0 {
		http.Error(w, "Ward metrics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics[0])
}

// GetVotingAllies returns voting alignment data for an official
func GetVotingAllies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	officialID := vars["id"]

	// Query votes from the official
	var officialVotes []map[string]interface{}
	_, err := db.Client.From("votes").
		Select("matter_id, vote_result", "exact", false).
		Eq("person_id", officialID).
		ExecuteTo(&officialVotes)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a map of matter_id -> vote_result for this official
	officialVoteMap := make(map[string]string)
	for _, vote := range officialVotes {
		if matterID, ok := vote["matter_id"].(string); ok {
			if voteResult, ok := vote["vote_result"].(string); ok {
				officialVoteMap[matterID] = voteResult
			}
		}
	}

	// Get all other officials
	var officials []models.Official
	_, err = db.Client.From("officials").
		Select("id, name, ward, party", "exact", false).
		ExecuteTo(&officials)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate alignment with each official
	type AllyData struct {
		OfficialID int     `json:"official_id"`
		Name       string  `json:"name"`
		Ward       *int    `json:"ward"`
		Party      string  `json:"party"`
		Alignment  float64 `json:"alignment"`
		Bloc       string  `json:"bloc"`
	}

	var allies []AllyData

	for _, other := range officials {
		if strconv.Itoa(other.ID) == officialID {
			continue
		}

		// Get votes for this official
		var otherVotes []map[string]interface{}
		_, err := db.Client.From("votes").
			Select("matter_id, vote_result", "exact", false).
			Eq("person_id", strconv.Itoa(other.ID)).
			ExecuteTo(&otherVotes)
		
		if err != nil {
			continue
		}

		// Calculate alignment
		matches := 0
		total := 0
		for _, vote := range otherVotes {
			if matterID, ok := vote["matter_id"].(string); ok {
				if officialVote, exists := officialVoteMap[matterID]; exists {
					if otherVote, ok := vote["vote_result"].(string); ok {
						total++
						if officialVote == otherVote {
							matches++
						}
					}
				}
			}
		}

		if total > 0 {
			alignment := (float64(matches) / float64(total)) * 100
			bloc := "Independent"
			if other.Party == "Democratic" {
				bloc = "Progressive Caucus"
			}

			allies = append(allies, AllyData{
				OfficialID: other.ID,
				Name:       other.Name,
				Ward:       other.Ward,
				Party:      other.Party,
				Alignment:  alignment,
				Bloc:       bloc,
			})
		}
	}

	// Sort by alignment (descending) and take top 10
	// Simple bubble sort for now
	for i := 0; i < len(allies); i++ {
		for j := i + 1; j < len(allies); j++ {
			if allies[j].Alignment > allies[i].Alignment {
				allies[i], allies[j] = allies[j], allies[i]
			}
		}
	}

	if len(allies) > 10 {
		allies = allies[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allies)
}

// GetRecentVotes returns recent votes for an official
func GetRecentVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	officialID := vars["id"]

	var votes []map[string]interface{}
	
	_, err := db.Client.From("votes").
		Select("*, matters(matter_name, matter_type)", "exact", false).
		Eq("person_id", officialID).
		Order("created_at", nil).
		Limit(10, "").
		ExecuteTo(&votes)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}
