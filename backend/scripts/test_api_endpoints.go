package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	baseURL := "https://webapi.legistar.com/v1/chicago"
	
	// Test different endpoints to find office/position information
	endpoints := []string{
		"/officeholders",
		"/officerecords",
		"/persons",
		"/bodies",
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	
	for _, endpoint := range endpoints {
		url := baseURL + endpoint
		fmt.Printf("\n=== Testing: %s ===\n", url)
		
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if resp.StatusCode != 200 {
			fmt.Printf("❌ Status: %d\n", resp.StatusCode)
			continue
		}
		
		// Try to parse and show first item
		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Printf("❌ Parse error: %v\n", err)
			continue
		}
		
		fmt.Printf("✅ Found %d items\n", len(data))
		if len(data) > 0 {
			fmt.Println("First item fields:")
			prettyJSON, _ := json.MarshalIndent(data[0], "  ", "  ")
			fmt.Printf("  %s\n", string(prettyJSON))
		}
	}
}
