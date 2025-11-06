package cityapi

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	
	if client.BaseURL != BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", BaseURL, client.BaseURL)
	}
	
	if client.HTTPClient == nil {
		t.Fatal("HTTPClient is nil")
	}
}

// Integration test - only run with -integration flag
func TestGetBodies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	client := NewClient()
	bodies, err := client.GetBodies()
	
	if err != nil {
		t.Logf("Warning: City API request failed: %v", err)
		t.Skip("City API may be unavailable")
	}
	
	if len(bodies) == 0 {
		t.Log("Warning: No bodies returned from API")
	} else {
		t.Logf("Successfully fetched %d bodies", len(bodies))
		
		// Check first body has expected fields
		if bodies[0].BodyName == "" {
			t.Error("First body has empty BodyName")
		}
	}
}

func TestGetMatters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	client := NewClient()
	
	// Fetch just one matter for testing
	matters, err := client.GetMatters(map[string]string{
		"$top": "1",
	})
	
	if err != nil {
		t.Logf("Warning: City API request failed: %v", err)
		t.Skip("City API may be unavailable")
	}
	
	if len(matters) > 0 {
		t.Logf("Successfully fetched matter: %s - %s", matters[0].MatterFile, matters[0].MatterTitle)
	}
}
