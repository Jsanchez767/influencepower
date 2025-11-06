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
	baseURL := "https://webapi.legistar.com/v1/chicago"
	
	// Test different queries on persons endpoint
	queries := []struct{
		name string
		params map[string]string
	}{
		{
			name: "All persons",
			params: map[string]string{},
		},
		{
			name: "Search: Gutierrez",
			params: map[string]string{"$search": "Gutierrez"},
		},
		{
			name: "Filter: PersonEmail contains Ward",
			params: map[string]string{"$filter": "contains(PersonEmail, 'Ward')"},
		},
		{
			name: "Filter: PersonEmail contains Ward14",
			params: map[string]string{"$filter": "contains(PersonEmail, 'Ward14')"},
		},
	}
	
	for _, query := range queries {
		log.Printf("\n📡 Testing: %s", query.name)
		
		// Build URL with query params
		queryParams := url.Values{}
		for key, value := range query.params {
			queryParams.Add(key, value)
		}
		
		endpoint := baseURL + "/persons"
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
			log.Printf("  Response: %s", string(body)[:200])
			continue
		}
		
		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("  ❌ JSON Error: %v", err)
			continue
		}
		
		log.Printf("  ✅ Status: %d, Count: %d", resp.StatusCode, len(data))
		
		if len(data) > 0 {
			log.Println("  📄 First 3 records:")
			for i := 0; i < min(3, len(data)); i++ {
				sampleJSON, _ := json.MarshalIndent(data[i], "    ", "  ")
				fmt.Println("    " + string(sampleJSON))
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
