package handlers

//*
// Test functions will test the followiing:
// 1. POST request to register a new webhook
// 2. GET request to get all webhooks
// 3. GET request to get a specific webhook
// 4. DELETE request to delete a specific webhook
//*


import (
	"assignment2/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var response utils.WebhookId

// TestPostRequest tests the POST request to register a new webhook
func TestPostRequest(t *testing.T) {
	// Create the payload
	payloadStruct := utils.RegisterWebhook{Url: "http://localhost:8081/invoked", Country: "SE", Event: "INVOKE"}
	jsonPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		t.Errorf("Error marshalling payload: %v", err)
	}

	// Create the request
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080"+utils.REGISTRATION_PATH, bytes.NewBuffer(jsonPayload))
	// Create the recorder
	w := httptest.NewRecorder()

	// Call the function with the testing flag
	registerNewWebhook(w, req, true)

	resp := w.Result()
	defer resp.Body.Close()

	// Test if the status code is what is expected
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusCreated)
		t.Errorf("Response body: %v", resp.Body)
	}
	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	// Test if the ID is not empty
	if response.ID == "" {
		t.Errorf("Response ID is empty")
	}
}

func TestPatchRequest(t *testing.T) {
	// payload to be changed from SE to NO
	payloadStruct := utils.RegisterWebhook{Url: "", Country: "NO", Event: ""}
	jsonPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		t.Errorf("Error marshalling payload: %v", err)
	}

	// make request
	req := httptest.NewRequest(http.MethodPatch, "http://localhost:8080"+utils.NOTIFICATION_PATH+response.ID, bytes.NewBuffer(jsonPayload))
	// recorder
	w := httptest.NewRecorder()

	// call function, get result and make sure its closed
	patchWebhook(w, req, true)

	resp := w.Result()
	defer resp.Body.Close()

	// function only returns 204 No Content
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusNoContent)
	}
}

// Test to get the webhook just created
func TestOneGetRequest(t *testing.T) {
	// Test the document just created
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080"+utils.NOTIFICATION_PATH+response.ID, nil)
	// Get the recorder
	w := httptest.NewRecorder()

	// Call the function with the flag
	getWebhooks(w, req, true)

	resp := w.Result()
	defer resp.Body.Close()

	// Check the response code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	// Get the response
	var response2 utils.ReturnWebhook
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &response2); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	// Check the response
	if response2.ID != response.ID {
		t.Errorf("Response ID was %v, expected %v", response2.ID, response.ID)
	}
}

// Test to get all webhooks registerd
func TestGetRequest(t *testing.T) {
	// Make the request
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080"+utils.NOTIFICATION_PATH, nil)
	// Get the recorder
	w := httptest.NewRecorder()

	// Call the function
	getWebhooks(w, req, true)

	// Result of the call
	resp := w.Result()
	defer resp.Body.Close()

	// Test if the status code is correct
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	// get the response body
	var response2 []utils.ReturnWebhook
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &response2); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	// Test if the response is not empty
	if len(response2) < 1 {
		t.Errorf("Response body was empty")
	}
}

// Delete the created document
func TestDeleteRequest(t *testing.T) {
	//make request
	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080"+utils.NOTIFICATION_PATH+response.ID, nil)
	//register the recorder
	w := httptest.NewRecorder()

	//call the function with the testing flag
	deleteWebhook(w, req, true)

	//get the response
	resp := w.Result()
	defer resp.Body.Close()

	//make sure the status code is what is expected
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status code was %v, expected %v", resp.StatusCode, http.StatusNoContent)
	}
}
