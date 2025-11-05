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
