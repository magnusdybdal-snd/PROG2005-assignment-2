package tests

import (
	"assignment2/utils"
	"encoding/json"
	"net/http"
	"testing"
)

const TESTSTRING = "d0rMwJGizFNrfmMPuhGK"

func TestAllGetRequest(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/dashboard/v1/notifications")
	if err != nil {
		t.Errorf("Error getting response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	var response []utils.ReturnWebhook
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	ids := make(map[string]bool)
	for _, v := range response {
		if v.ID == "" {
			t.Errorf("Response ID is empty")
		}
		if len(v.ID) != 20 {
			t.Errorf("Response ID is not 20 characters long")
		}
		if ids[v.ID] {
			t.Errorf("Response ID %v is not unique", v.ID)
		}
		ids[v.ID] = true
	}
}

func TestIdGetRequest(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/dashboard/v1/notifications/" + TESTSTRING)
	if err != nil {
		t.Errorf("Error getting response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	var response utils.ReturnWebhook
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if response.ID != TESTSTRING {
		t.Errorf("Response ID was %v, expected %v", response.ID, TESTSTRING)
	}
}
