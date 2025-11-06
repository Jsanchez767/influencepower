package api

import (
	"encoding/json"
	"net/http"

	"github.com/Jsanchez767/InfluencePower/backend/handlers"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *mux.Router) {
	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", HealthCheck).Methods("GET")

	// Officials routes
	api.HandleFunc("/officials", handlers.GetOfficials).Methods("GET")
	api.HandleFunc("/officials", handlers.CreateOfficial).Methods("POST")
	api.HandleFunc("/officials/{id}", handlers.GetOfficialByID).Methods("GET")
	api.HandleFunc("/officials/{id}", handlers.UpdateOfficial).Methods("PUT")
	api.HandleFunc("/officials/{id}", handlers.DeleteOfficial).Methods("DELETE")
	api.HandleFunc("/officials/party/{party}", handlers.GetOfficialsByParty).Methods("GET")
	api.HandleFunc("/officials/ward/{ward}", handlers.GetOfficialsByWard).Methods("GET")

	// Voting records routes
	api.HandleFunc("/officials/{id}/voting-records", handlers.GetVotingRecords).Methods("GET")
	api.HandleFunc("/voting-records", handlers.CreateVotingRecord).Methods("POST")

	// Ward statistics routes
	api.HandleFunc("/wards/{ward}/statistics", handlers.GetWardStatistics).Methods("GET")

	// Committees routes
	api.HandleFunc("/committees", handlers.GetCommittees).Methods("GET")
	api.HandleFunc("/officials/{id}/committees", handlers.GetOfficialCommittees).Methods("GET")

	// Metrics routes
	api.HandleFunc("/officials/{id}/metrics", handlers.GetOfficialMetrics).Methods("GET")
	api.HandleFunc("/wards/{ward}/metrics", handlers.GetWardMetrics).Methods("GET")
	api.HandleFunc("/officials/{id}/voting-allies", handlers.GetVotingAllies).Methods("GET")
	api.HandleFunc("/officials/{id}/recent-votes", handlers.GetRecentVotes).Methods("GET")
}

// HealthCheck returns the health status of the API
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"message": "InfluencePower API is running",
	})
}
