package handlers

import (
	"assignment2/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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

func TestGetRestCountriesData(t *testing.T) {
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
}
