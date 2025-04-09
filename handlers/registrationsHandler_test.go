package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestUpdateDocument(t *testing.T) {

}
func TestDeleteDocument(t *testing.T) {

}

func TestRegisterDashConfig(t *testing.T) {
	// --- Test Setup ---
	// Disable log output during tests
	originalOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(originalOutput)

	// Initialize Firestore - This should connect to your TEST environment/emulator
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore for test: %v", err)
	}
	defer utils.CloseFirestore()

	// Test Cases
	testCases := []struct {
		name             string
		inputBody        string   // Request body JSON as string
		wantStatusCode   int      // Expected HTTP status code
		wantBodyContains []string // Substrings expected in the response body
		checkFirestore   bool     // Should we verify Firestore write on success?
		expectedCountry  string   // Expected country value for Firestore check
		expectedIsoCode  string   // Expected isoCode value for Firestore check
	}{
		{
			name: "Success - Valid Payload",
			// JSON structure must match utils.DashboardAlteration for validation
			inputBody: `{
                "country": "Test Success",
                "isoCode": "TS",
                "features": {
                    "temperature": true,
                    "precipitation": true,
                    "capital": false,
                    "coordinates": true,
                    "population": false,
                    "area": true,
                    "targetCurrencies": ["NOK"]
                }
            }`,
			wantStatusCode:   http.StatusCreated,                 // Handler returns 201 on success
			wantBodyContains: []string{`"id":`, `"lastChange":`}, // Response contains ID and timestamp
			checkFirestore:   true,                               // Verify Firestore write
			expectedCountry:  "Test Success",                     // Data expected in Firestore
			expectedIsoCode:  "TS",
		},
		{
			name:           "Fail - Invalid JSON Syntax",
			inputBody:      `{"country": "Test Fail", "isoCode": "TF",`, // Malformed JSON
			wantStatusCode: http.StatusBadRequest,                       // Caught by ValidatePostRequest
			// Expect error message from ValidatePostRequest
			wantBodyContains: []string{"Unknown fields are present"},
			checkFirestore:   false,
		},
		{
			name: "Fail - Unknown Field",
			inputBody: `{
                "country": "Test Fail Unknown",
                "isoCode": "TU",
                "features": {},
                "extraField": 123
            }`, // Contains a field not in DashboardAlteration
			wantStatusCode:   http.StatusBadRequest,      // Caught by ValidatePostRequest
			wantBodyContains: []string{"Unknown fields"}, // From ValidatePostRequest
			checkFirestore:   false,
		},
		{
			name:             "Fail - Empty Body",
			inputBody:        ``,
			wantStatusCode:   http.StatusBadRequest,      // Caught by ValidatePostRequest (likely EOF)
			wantBodyContains: []string{"Unknown fields"}, // From ValidatePostRequest
			checkFirestore:   false,
		},
		{
			name: "Fail - Empty JSON Object",
			// ValidatePostRequest might pass this if fields are optional,
			// but let's assume 'country' or similar is implicitly required or fails later.
			inputBody:        `{}`,
			wantStatusCode:   http.StatusBadRequest,               // Assuming validation requires some fields
			wantBodyContains: []string{"Missing required fields"}, // Or maybe a different error later if validation passes
			checkFirestore:   false,
		},
		{
			name: "Fail - Type Mismatch",
			inputBody: `{
				"country": "Test Type",
                "isoCode": "TT",
				"features": {"temperature": "true"}
			}`, // temperature should be bool
			wantStatusCode:   http.StatusBadRequest,      // Caught by ValidatePostRequest
			wantBodyContains: []string{"Unknown fields"}, // From ValidatePostRequest
			checkFirestore:   false,
		},
	}

	// Loops over test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/registrations", strings.NewReader(tc.inputBody)) // Use actual endpoint path
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctx := context.Background()

			// Call the handler function being tested
			registerDashConfig(w, req, ctx, true)

			// Check Status Code
			if w.Code != tc.wantStatusCode {
				t.Errorf("want status code %d, got %d. Body: %s", tc.wantStatusCode, w.Code, w.Body.String())
			}

			// Check Response Body Substrings
			bodyStr := w.Body.String()
			for _, sub := range tc.wantBodyContains {
				if !strings.Contains(bodyStr, sub) {
					t.Errorf("want body to contain %q, got %q", sub, bodyStr)
				}
			}

			// Optional: Check Firestore Write on Success
			if tc.checkFirestore && w.Code == http.StatusCreated {
				var respData utils.RegistrationGetResponse
				if err := json.Unmarshal(w.Body.Bytes(), &respData); err != nil {
					t.Fatalf("Failed to unmarshal success response body %q: %v", bodyStr, err)
				}
				if respData.Id == "" {
					t.Fatalf("Success response body did not contain an 'id': %s", bodyStr)
				}
				if respData.LastRetrieval.IsZero() {
					t.Fatalf("Success response body did not contain a valid 'lastChange': %s", bodyStr)
				}

				// Assumes utils.DASHBOARD_COLLECTION is set to the TEST collection name
				docRef := utils.FirestoreClient.Collection(utils.DASHBOARD_TEST_COLLECTION).Doc(respData.Id)
				docSnap, err := docRef.Get(ctx)
				// Clean up the document afterwards, regardless of test outcome (if possible)
				defer func() {
					if _, err := docRef.Delete(ctx); err != nil {
						t.Logf("Warning: Failed to delete test document %s: %v", respData.Id, err)
					}
				}() // End defer

				if err != nil {
					t.Fatalf("Firestore check: Failed to get document with ID %s: %v", respData.Id, err)
				}
				if !docSnap.Exists() {
					t.Fatalf("Firestore check: Document with ID %s was not found", respData.Id)
				}

				// Compare Firestore data
				var firestoreData utils.DashboardConfig
				if err := docSnap.DataTo(&firestoreData); err != nil {
					t.Fatalf("Firestore check: Failed to read data into struct: %v", err)
				}

				if firestoreData.Country != tc.expectedCountry {
					t.Errorf("Firestore check: Country mismatch, got %q, want %q", firestoreData.Country, tc.expectedCountry)
				}
				if firestoreData.IsoCode != tc.expectedIsoCode {
					t.Errorf("Firestore check: IsoCode mismatch, got %q, want %q", firestoreData.IsoCode, tc.expectedIsoCode)
				}
				// Check timestamp is recent (within ~5 seconds of the response timestamp)
				if firestoreData.LastRetrieval.IsZero() || firestoreData.LastRetrieval.Unix()-respData.LastRetrieval.Unix() > 5 || respData.LastRetrieval.Unix()-firestoreData.LastRetrieval.Unix() > 5 {
					t.Errorf("Firestore check: LastRetrieval timestamp mismatch or too different. Got %v, Response was %v", firestoreData.LastRetrieval, respData.LastRetrieval)
				}
			}
		})
	}
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
