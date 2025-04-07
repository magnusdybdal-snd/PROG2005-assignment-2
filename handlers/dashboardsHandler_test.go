package handlers

import (
	"assignment2/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Uses test tables to test the function calculateMean(...)
func TestCalculateMean(t *testing.T) {
	// Defines test cases for calculateMean(...)
	testCases := []struct {
		name  string    // Name of the test case
		input []float64 // Input slice for calculation
		want  float64   // Expected result
	}{
		{name: "Empty Slice", input: []float64{}, want: 0.0},
		{name: "Single Element", input: []float64{5.5}, want: 5.50},
		{name: "Multiple Positive Integers", input: []float64{1.0, 2.0, 3.0}, want: 2.00},
		{name: "Mixed Positive and Negative", input: []float64{-1.0, 1.0, 3.0, 5.0}, want: 2.00},
		{name: "Needs Rounding Up", input: []float64{1.111, 2.222, 3.333}, want: 2.22},
		{name: "Needs Rounding Up (midpoint)", input: []float64{1.115, 2.225}, want: 1.67},
		{name: "Needs Rounding Down", input: []float64{10.0, 20.0, 35.0}, want: 21.67},
		{name: "Already Two Decimal", input: []float64{1.23, 4.56, 7.89}, want: 4.56},
		{name: "Zeros", input: []float64{0.0, 0.0, 0.0}, want: 0.00},
	}

	// Iterates over the test cases in the test table above
	for _, tc := range testCases {

		t.Run(tc.name, func(st *testing.T) {
			got := calculateMean(tc.input) // Run calculation for each testcase
			if got != tc.want {            // Check if the result matches the expected value
				t.Errorf("calculateMean(%v) == %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

/*
*	Function with tests for getCountriesData()
 */
func TestGetRestCountriesData(t *testing.T) {
	// --- API returns 200 OK and correct response
	// Uses t.Run for more customasation in go test
	t.Run("Success path for Norway", func(t *testing.T) {
		isoCode := "no"
		// Setting up tests with mocked data.
		jsonPath := filepath.Join("testdata", "RestCountriesNorwaySuccess.json")
		mockJSONResponse, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("TEST SETUP FAILED: Could not read testfile %s: %v", jsonPath, err)
		}
		// Initialises a new server for testing and writes the response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(mockJSONResponse)
		}))
		defer server.Close()
		// The expected response from CountriesNow (no)
		expectedData := utils.RestCountriesResponse{
			Capital:     []string{"Oslo"},
			Coordinates: []int{62, 10},
			Population:  5379475,
			Area:        323802,
			Currencies: map[string]interface{}{
				"NOK": map[string]interface{}{
					"name":   "Norwegian krone",
					"symbol": "kr",
				},
			},
		}

		testClient := server.Client()
		actualData, actualErr := getRestCountriesData(testClient, server.URL+"/v3.1/alpha/", isoCode)

		if actualErr != nil {
			t.Fatalf("getRestCountriesData() returned an unexpected error: %v", actualErr)
		}
		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("getRestCountriesData() returned unexpected data.\nGot:\n%#v\nWant:\n%#v", actualData, expectedData)
		}
	})

	// --- API returns 404 Not Found Scenario
	// Uses t.Run for more customasation in go test
	t.Run("API returns 404", func(t *testing.T) {
		isoCode := "zz" // Invalid iso code to cause 404

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound) // 404
		}))
		defer server.Close()

		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+"/v3.1/alpha/", isoCode)

		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for 404 status, but got nil")
		}
		expectedErrorMsg := fmt.Sprintf("API returned non-200 status code: %d", http.StatusNotFound)
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// --- Malformed JSON Response Scenario
	t.Run("Malformed JSON response", func(t *testing.T) {
		isoCode := "no" // Iso code does not matter as we will write malformed JSON directly

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"name": {"common": "Norway"}, "capital": ["Oslo"]`))
		}))
		defer server.Close()

		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+"/v3.1/alpha/", isoCode)

		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for malformed JSON but got nil")
		}
		expectedErrorMsg := "error when decoding json"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// -- Empty JSON Array Response Scenario
	t.Run("Empty JSON array response", func(t *testing.T) {
		isoCode := "no"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`)) // Empty Json array
		}))
		defer server.Close()

		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+"/v3.1/alpha/", isoCode)

		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for empty JSON array, but got nil")
		}

		expectedErrorMsg := "apiResponse slice is empty"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// -- Network Error Scenario
	t.Run("Network Error", func(t *testing.T) {
		isoCode := "no"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This handler should never be called during this test case
			t.Errorf("UNEXPECTED: Network error test server recieved a request: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		// Get the URL for testing before closing the server
		closedServerURL := server.URL + "/v3.1/alpha/"
		server.Close()

		// Create a standard http client. Not using server.Client() here as we want a to attempt a network connection
		testClient := &http.Client{}
		_, actualErr := getRestCountriesData(testClient, closedServerURL, isoCode)

		if actualErr == nil {
			t.Fatalf("getRestCountriesData() expected a network error when connecting to %s, but got nil", closedServerURL)
		}

		expectedErrorSubstrings := []string{"connection refused", "connect: connection refused"}
		errorMatched := false
		errStr := actualErr.Error() // Get the error message string

		for _, sub := range expectedErrorSubstrings {
			if strings.Contains(errStr, sub) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			// Use t.Errorf for assertion failures.
			t.Errorf("getRestCountriesData() error = %q, did not contain expected network error substrings (%v)", actualErr, expectedErrorSubstrings)
		}
	})
}
