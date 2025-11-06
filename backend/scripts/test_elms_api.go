package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	baseURL := "https://api.chicityclerkelms.chicago.gov"
	
	log.Println("🔍 Testing Chicago City Clerk ELMS API")
	log.Println("   Base URL:", baseURL)
	log.Println()
	
	// Test the export/person endpoint with ward filter
	tests := []struct {
		name   string
		endpoint string
		params map[string]string
	}{
		{
			name:     "Ward 14 filter",
			endpoint: "/export/person",
			params:   map[string]string{"filter": "ward eq '14'"},
		},
		{
			name:     "All persons",
			endpoint: "/export/person",
			params:   map[string]string{},
		},
		{
			name:     "Search Gutierrez",
			endpoint: "/export/person",
			params:   map[string]string{"search": "Gutierrez"},
		},
	}
	
	for _, test := range tests {
		log.Printf("\n📡 Test: %s", test.name)
		
		queryParams := url.Values{}
		for key, value := range test.params {
			queryParams.Add(key, value)
		}
		
		endpoint := baseURL + test.endpoint
		if len(queryParams) > 0 {
			endpoint = endpoint + "?" + queryParams.Encode()
		}
		
		log.Printf("   URL: %s", endpoint)
		
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Printf("  ❌ Error: %v", err)
			continue
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		
		if resp.StatusCode != 200 {
			log.Printf("  ❌ Status: %d", resp.StatusCode)
			if len(body) < 500 {
				log.Printf("  Response: %s", string(body))
			}
			continue
		}
		
		// Try to parse as JSON
		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("  ❌ JSON Error: %v", err)
			log.Printf("  Response preview (first 500 chars): %s", string(body)[:min(500, len(body))])
			continue
		}
		
		log.Printf("  ✅ SUCCESS! Status: %d, Count: %d", resp.StatusCode, len(data))
		
		if len(data) > 0 {
			log.Println("\n  📄 First record:")
			sampleJSON, _ := json.MarshalIndent(data[0], "    ", "  ")
			fmt.Println("    " + string(sampleJSON))
			
			log.Println("\n  📋 Available fields:")
			for key := range data[0] {
				fmt.Printf("    - %s\n", key)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
