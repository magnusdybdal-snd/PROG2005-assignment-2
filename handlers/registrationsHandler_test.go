package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestHandleMessages(t *testing.T) {

}

func TestUpdateDocument(t *testing.T) {

}
func TestDeleteDocument(t *testing.T) {

}
func TestRegisterDashConfig(t *testing.T) {

}

func TestDisplayDocument(t *testing.T) {

	testCases := []struct {
		name       string
		requestURL string
		pathID     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Get single document",
			requestURL: utils.REGISTRATION_PATH + "Qqx0bRWYvE6J3mOhDEMF",
			pathID:     "/Users/maseilertsen/Documents/04_code-projects/Golang/assignment2/stud/testdata/documents/displayResponse.json", // Absolute path to control-fil.
			wantStatus: http.StatusOK,
		},
		{
			name:       "Get all document",
			requestURL: utils.REGISTRATION_PATH,
			pathID:     "/Users/maseilertsen/Documents/04_code-projects/Golang/assignment2/stud/testdata/documents/allDocuments.json", // Absolute path to control-fil.
			wantStatus: http.StatusOK,
		},
	}

	// Initialize Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer utils.CloseFirestore()

	// Run all test-scenarios
	for _, tc := range testCases {
		// Create request.
		req := httptest.NewRequest(http.MethodGet, tc.requestURL, nil)

		// Create response recorder.
		w := httptest.NewRecorder()

		// Create a context
		ctx := context.Background()

		// Calling the function
		displayDocument(w, req, ctx, true)

		if got, want := w.Code, tc.wantStatus; got != want {
			t.Errorf("Status code: got %v, want %v", got, want)
		}

		expected, err := os.ReadFile(tc.pathID)
		if err != nil {
			t.Fatalf("Failed to read expected response file: %v", err)
		}

		if tc.wantStatus == http.StatusOK {
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &gotJSON); err != nil {
				t.Fatalf("failed to parse response JSON: %v", err)
			}
			if err := json.Unmarshal(expected, &wantJSON); err != nil {
				t.Fatalf("failed to parse expected JSON: %v", err)
			}

			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("JSON response doesn't match expected")
			}

		}
	}
}
