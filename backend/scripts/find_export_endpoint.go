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
	// Test different base URLs for the export endpoint
	testURLs := []string{
		"https://webapi.legistar.com/v1/chicago/export/person",
		"https://chicago.legistar.com/export/person",
		"https://webapi.legistar.com/export/person",
		"https://chicago.councilmatic.org/export/person",
	}
	
	filter := "ward eq '14'"
	
	for _, baseURL := range testURLs {
		log.Printf("\n📡 Testing: %s", baseURL)
		log.Printf("   Filter: %s", filter)
		
		// Build URL with filter
		queryParams := url.Values{}
		queryParams.Add("filter", filter)
		
		endpoint := baseURL + "?" + queryParams.Encode()
		log.Printf("   Full URL: %s", endpoint)
		
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
		
		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("  ❌ JSON Error: %v", err)
			log.Printf("  Response preview: %s", string(body)[:min(200, len(body))])
			continue
		}
		
		log.Printf("  ✅ SUCCESS! Status: %d, Count: %d", resp.StatusCode, len(data))
		
		if len(data) > 0 {
			log.Println("  📄 Sample record:")
			sampleJSON, _ := json.MarshalIndent(data[0], "    ", "  ")
			fmt.Println("    " + string(sampleJSON))
		}
		
		return // Found it!
	}
	
	log.Println("\n❌ Could not find working export endpoint")
	log.Println("💡 Can you share the exact URL that worked for you?")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
