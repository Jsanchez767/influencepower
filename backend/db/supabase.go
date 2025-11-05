package db

import (
	"log"

	postgrest "github.com/supabase-community/postgrest-go"
)

var Client *postgrest.Client

// InitSupabase initializes the Supabase client
func InitSupabase(url, key string) {
	Client = postgrest.NewClient(url+"/rest/v1", "", map[string]string{
		"apikey":        key,
		"Authorization": "Bearer " + key,
	})
	
	if Client == nil {
		log.Fatal("Failed to initialize Supabase client")
	}
	
	log.Println("Supabase client initialized successfully")
}
