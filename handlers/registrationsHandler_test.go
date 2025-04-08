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


func TestUpdateDocument(t *testing.T) {

}
func TestDeleteDocument(t *testing.T) {

}
func TestRegisterDashConfig(t *testing.T) {

		// Initialize Firestore
		if err := utils.InitFirestore(); err != nil {
			log.Fatalf("Error initializing Firestore: %v", err)
		}
		defer utils.CloseFirestore()
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
			requestURL: utils.REGISTRATION_PATH,                                                                                       // No message string to trigger iteration.
			pathID:     "/Users/maseilertsen/Documents/04_code-projects/Golang/assignment2/stud/testdata/documents/allDocuments.json", // Absolute path to control-fil.
			wantStatus: http.StatusOK,
		},
		{
			name:       "Get non-existent document",
			requestURL: utils.REGISTRATION_PATH + "this-id-does-not-exist-123", // Use an ID guaranteed not to exist
			// pathID: "", // No JSON file expected for error
			wantStatus: http.StatusBadRequest,                                     // Expect 400
			wantBody:   http.StatusText(http.StatusBadRequest) + ": Invalid ID\n", // Expected error message from http.Error (it adds a newline)
		},
	}

	// Initialize Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer utils.CloseFirestore()

	// Run all test-scenarios
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { // Use t.Run for better sub-test naming
			// Create request.
			req := httptest.NewRequest(http.MethodGet, tc.requestURL, nil)

			// Create response recorder.
			w := httptest.NewRecorder()

			// Create a context
			ctx := context.Background()

			// Calling the function
			displayDocument(w, req, ctx, true) // Pass test=true

			// Check Status Code (keep this as is)
			if got, want := w.Code, tc.wantStatus; got != want {
				t.Errorf("Status code: got %v, want %v", got, want)
			}
			// Check if wanted status is OK (200)
			if tc.wantStatus == http.StatusOK {
				// Check JSON body for OK status
				if tc.pathID == "" {
					t.Fatal("pathID (expected JSON file) must be set for OK test cases")
				}
				expected, err := os.ReadFile(tc.pathID)
				if err != nil {
					t.Fatalf("Failed to read expected response file '%s': %v", tc.pathID, err)
				}

				var gotJSON, wantJSON interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &gotJSON); err != nil {
					// Log the body content if unmarshalling fails
					t.Fatalf("failed to parse actual response JSON: %v\nBody: %s", err, w.Body.String())
				}
				if err := json.Unmarshal(expected, &wantJSON); err != nil {
					t.Fatalf("failed to parse expected JSON from file '%s': %v", tc.pathID, err)
				}

				log.Printf("Test %s: got JSON: %v", tc.name, gotJSON) // Use log.Printf with test name
				log.Printf("Test %s: want JSON: %v", tc.name, wantJSON)

				if !reflect.DeepEqual(gotJSON, wantJSON) {
					// Consider using a diff library here for better output on complex JSON mismatches
					t.Errorf("JSON response doesn't match expected")
				}

			} else {
				// Check Plain Text body for non-OK status (errors)
				if tc.wantBody != "" { // Only check if an expected error body is defined
					if gotBody := w.Body.String(); gotBody != tc.wantBody {
						t.Errorf("Error response body: got %q, want %q", gotBody, tc.wantBody)
					}
				}
				// If tc.wantBody is empty for an error case, we don't check the body content.
			}
		})
	} // End of loop
}
