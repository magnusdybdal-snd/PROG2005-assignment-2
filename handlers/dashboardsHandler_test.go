package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("FATAL: TestMain failed to initialize Firestore via utils.InitFirestore: %v", err)
	}
	log.Println("Firestore client initialized for utils tests")

	exitCode := m.Run()

	utils.CloseFirestore()
	log.Println("Firestore client closed.")
	os.Exit(exitCode)
}

// Tests uses t.Run for naming tests and more customizing when running test

/*
*	Function with tests for HandleGetDashboard()
 */
func TestHandleGetDashboard(t *testing.T) {
	// === Test Case 1: Success with all features enabled ===
	t.Run("Success with all features", func(t *testing.T) {

		// 1. Define mock data and expected results
		testID := "test-id-all-features"
		mockConfig := utils.DashboardConfig{
			Country: "Testland", IsoCode: "TL",
			Features: struct {
				Temperature      bool     "firestore:\"temperature\" json:\"temperature\""
				Precipitation    bool     "firestore:\"precipitation\" json:\"precipitation\""
				Capital          bool     "firestore:\"capital\" json:\"capital\""
				Coordinates      bool     "firestore:\"coordinates\" json:\"coordinates\""
				Population       bool     "firestore:\"population\" json:\"population\""
				Area             bool     "firestore:\"area\" json:\"area\""
				TargetCurrencies []string "firestore:\"targetCurrencies\" json:\"targetCurrencies\""
			}{
				Temperature: true, Precipitation: true, Capital: true, Coordinates: true, Population: true,
				Area: true, TargetCurrencies: []string{"USD", "EUR", "SEK"},
			},
		}
		mockCountriesData := utils.RestCountriesResponse{
			Capital: []string{"Testville"}, Coordinates: []float64{20.0, 30.0}, Population: 1234567, Area: 10000,
			Currencies: map[string]interface{}{"TLD": map[string]interface{}{"name": "Test Dollar"}},
		}
		mockMetroData := utils.MetroMeanValues{MeanTemperature: 15.5, MeanPrecipitation: 2.3}
		mockCurrencyData := map[string]float64{"USD": 1.1, "EUR": 0.9, "SEK": 8.3}

		expectedResponse := utils.DashboardResponse{
			Country: "Testland", IsoCode: "TL",
			Features: struct {
				Temperature      float64            "json:\"temperature,omitempty\""
				Precipitation    float64            "json:\"precipitation,omitempty\""
				Capital          string             "json:\"capital,omitempty\""
				Coordinates      map[string]float64 "json:\"coordinates,omitempty\""
				Population       int                "json:\"population,omitempty\""
				Area             float64            "json:\"area,omitempty\""
				TargetCurrencies map[string]float64 "json:\"targetCurrencies,omitempty\""
			}{
				Temperature: 15.5, Precipitation: 2.3, Capital: "Testville",
				Coordinates: map[string]float64{"latitude": 20.0, "longitude": 30.0},
				Population:  1234567, Area: 10000, TargetCurrencies: map[string]float64{"USD": 1.1, "EUR": 0.9, "SEK": 8.3},
			},
		}

		// 2. Setup mocks for the external dependent function calls
		originalGetConfig := getConficFunc
		origianlGetCountries := getCountriesFunc
		originalGetMetro := getMetroFunc
		origianlGetCurrency := getCurrencyFunc
		t.Cleanup(func() { // Restores the functions after testing
			getConficFunc = originalGetConfig
			getCountriesFunc = origianlGetCountries
			getMetroFunc = originalGetMetro
			getCurrencyFunc = origianlGetCurrency
		})

		// 3. Seting up the functions to return the mocked responses above
		getConficFunc = func(ctx context.Context, id string, collection string) (utils.DashboardConfig, error) {
			return mockConfig, nil
		}
		getCountriesFunc = func(client *http.Client, baseURL string, IsoCode string) (utils.RestCountriesResponse, error) {
			return mockCountriesData, nil
		}
		getMetroFunc = func(client *http.Client, baseURL string, lat, long float64) (utils.MetroMeanValues, error) {
			return mockMetroData, nil
		}
		getCurrencyFunc = func(client *http.Client, baseURL string, currencies map[string]interface{}, targetCurrencies []string) (map[string]float64, error) {
			return mockCurrencyData, nil
		}

		// 4. Setting up request/recorder
		req := httptest.NewRequest(http.MethodGet, utils.DASHBOARD_PATH+testID, nil)
		req.SetPathValue("id", testID)
		w := httptest.NewRecorder()

		// 5. Execute handler
		HandleGetDashboard(w, req)

		// 6. Assertions
		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, w.Body.String())
		}
		if ctype := w.Header().Get("Content-Type"); ctype != "application/json" {
			t.Errorf("handler returned wrong content type: got %qm want %q", ctype, "application/json")
		}

		var actualResponse utils.DashboardResponse
		if err := json.NewDecoder(w.Body).Decode(&actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v. Body: %s", err, w.Body.String())
		}
		// Zero out time for comparison
		actualResponse.LastRetrieval = time.Time{}
		expectedResponse.LastRetrieval = time.Time{}

		if !reflect.DeepEqual(actualResponse, expectedResponse) {
			t.Errorf("handler returned unexpected body:\nGot:\n%#v\nWant:%#v", actualResponse, expectedResponse)
		}
	})

	// === Test Case 2: Success with partial features (Temperature & Precipiation) ===
	t.Run("Success with temperature and precipiation", func(t *testing.T) {
		// 1. Define mock data and expected results
		testID := "test-id-temp-prec"
		mockConfig := utils.DashboardConfig{
			Country: "WeatherTemp", IsoCode: "WT",
			Features: struct {
				Temperature      bool     "firestore:\"temperature\" json:\"temperature\""
				Precipitation    bool     "firestore:\"precipitation\" json:\"precipitation\""
				Capital          bool     "firestore:\"capital\" json:\"capital\""
				Coordinates      bool     "firestore:\"coordinates\" json:\"coordinates\""
				Population       bool     "firestore:\"population\" json:\"population\""
				Area             bool     "firestore:\"area\" json:\"area\""
				TargetCurrencies []string "firestore:\"targetCurrencies\" json:\"targetCurrencies\""
			}{
				Temperature: true, Precipitation: true,
			},
		}
		mockCountriesData := utils.RestCountriesResponse{Coordinates: []float64{62.0, 10.0}}
		mockMetroData := utils.MetroMeanValues{MeanTemperature: 15.5, MeanPrecipitation: 2.3}
		expectedResponse := utils.DashboardResponse{
			Country: "WeatherTemp", IsoCode: "WT",
			Features: struct {
				Temperature      float64            "json:\"temperature,omitempty\""
				Precipitation    float64            "json:\"precipitation,omitempty\""
				Capital          string             "json:\"capital,omitempty\""
				Coordinates      map[string]float64 "json:\"coordinates,omitempty\""
				Population       int                "json:\"population,omitempty\""
				Area             float64            "json:\"area,omitempty\""
				TargetCurrencies map[string]float64 "json:\"targetCurrencies,omitempty\""
			}{
				Temperature: 15.5, Precipitation: 2.3,
			},
		}

		// Flags to check if unnecessary APIs were called
		countriesCalled := false
		metroCalled := false
		currencyCalled := false

		// 2. Setup mocks for the external dependent function calls
		originalGetConfig := getConficFunc
		origianlGetCountries := getCountriesFunc
		originalGetMetro := getMetroFunc
		origianlGetCurrency := getCurrencyFunc
		t.Cleanup(func() { // Restores the functions after testing
			getConficFunc = originalGetConfig
			getCountriesFunc = origianlGetCountries
			getMetroFunc = originalGetMetro
			getCurrencyFunc = origianlGetCurrency
		})

		// 3. Seting up the functions to return the mocked responses above
		getConficFunc = func(ctx context.Context, id string, collection string) (utils.DashboardConfig, error) {
			return mockConfig, nil
		}
		getCountriesFunc = func(client *http.Client, baseURL string, IsoCode string) (utils.RestCountriesResponse, error) {
			countriesCalled = true
			return mockCountriesData, nil
		}
		getMetroFunc = func(client *http.Client, baseURL string, lat, long float64) (utils.MetroMeanValues, error) {
			metroCalled = true
			return mockMetroData, nil
		}
		getCurrencyFunc = func(client *http.Client, baseURL string, currencies map[string]interface{}, targetCurrencies []string) (map[string]float64, error) {
			return nil, nil // Should not be called
		}

		// 4. Setting up request/recorder
		req := httptest.NewRequest(http.MethodGet, utils.DASHBOARD_PATH+testID, nil)
		req.SetPathValue("id", testID)
		w := httptest.NewRecorder()

		// 5. Execute handler
		HandleGetDashboard(w, req)

		// 6. Assertions
		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, w.Body.String())
		}

		if !countriesCalled {
			t.Error("getRestCountriesData was expected but not called")
		}
		if !metroCalled {
			t.Error("getMetroData was expected but not called")
		}
		if currencyCalled {
			t.Error("getCurrencyData was called unexpectedly")
		}

		if ctype := w.Header().Get("Content-Type"); ctype != "application/json" {
			t.Errorf("handler returned wrong content type: got %qm want %q", ctype, "application/json")
		}

		var actualResponse utils.DashboardResponse
		if err := json.NewDecoder(w.Body).Decode(&actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v. Body: %s", err, w.Body.String())
		}
		// Zero out time for comparison
		actualResponse.LastRetrieval = time.Time{}
		expectedResponse.LastRetrieval = time.Time{}

		if !reflect.DeepEqual(actualResponse, expectedResponse) {
			t.Errorf("handler returned unexpected body:\nGot:\n%#v\nWant:%#v", actualResponse, expectedResponse)
		}
	})

	// === Test Case 3: Error: No ID parameter ===
	t.Run("Error missing ID parameter", func(t *testing.T) {
		// 1. Setting up request/recorder without the id in path
		req := httptest.NewRequest(http.MethodGet, utils.DASHBOARD_PATH, nil)
		w := httptest.NewRecorder()

		// 2. Execute handler
		HandleGetDashboard(w, req)

		// 3. Assertions
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v, want %v", w.Code, http.StatusBadRequest)
		}
		expectedMessage := "Dashboard id is required."
		if body := w.Body.String(); !strings.Contains(body, expectedMessage) {
			t.Errorf("hanlder returned unexpected body: got %q want substring %q", body, expectedMessage)
		}
	})
}

/*
 *	Function with tests for getCountriesData()
 */
func TestGetRestCountriesData(t *testing.T) {
	// Define standard input(s) used across multiple tests
	defaultIsoCode := "no"
	APIString := "/v3.1/alpha/"

	// === Test Case 1: Success Path ===
	t.Run("Success country path for Norway", func(t *testing.T) {
		// 1. Setting up tests with mock JSON response.
		jsonPath := filepath.Join("testdata", "RestCountriesNorwaySuccess.json")
		mockJSONResponse, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("TEST SETUP FAILED: Could not read testfile %s: %v", jsonPath, err)
		}

		// 2. Set up test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(mockJSONResponse)
		}))
		defer server.Close()

		// 3. Define the expected response
		expectedData := utils.RestCountriesResponse{
			Capital:     []string{"Oslo"},
			Coordinates: []float64{62.0, 10.0},
			Population:  5379475,
			Area:        323802,
			Currencies: map[string]interface{}{
				"NOK": map[string]interface{}{
					"name":   "Norwegian krone",
					"symbol": "kr",
				},
			},
		}

		// 4. Call the function under test
		testClient := server.Client()
		actualData, actualErr := getRestCountriesData(testClient, server.URL+APIString, defaultIsoCode)

		// 5. Assertions
		if actualErr != nil {
			t.Fatalf("getRestCountriesData() returned an unexpected error: %v", actualErr)
		}
		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("getRestCountriesData() returned unexpected data.\nGot:\n%#v\nWant:\n%#v", actualData, expectedData)
		}
	})

	// === Test Case 2: API Returns 404 Status ==
	t.Run("API returns 404", func(t *testing.T) {
		isoCode := "zz" // Invalid iso code to cause 404

		// 1. Setup test server to return 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound) // 404
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+APIString, isoCode)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for 404 status, but got nil")
		}
		expectedErrorMsg := fmt.Sprintf("API returned non-200 status code: %d", http.StatusNotFound)
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 3: Malformed JSON Response ===
	t.Run("Malformed JSON response", func(t *testing.T) {
		// 1. Setup test server to return bad JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Invalid JSON
			_, _ = w.Write([]byte(`[{"name": {"common": "Norway"}, "capital": ["Oslo"]`))
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+APIString, defaultIsoCode)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for malformed JSON but got nil")
		}
		expectedErrorMsg := "error when decoding json"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 4: Empty JSON Array ===
	t.Run("Empty JSON array response", func(t *testing.T) {
		// 1. Setup test server with empty response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`)) // Empty Json array
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getRestCountriesData(testClient, server.URL+APIString, defaultIsoCode)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getRestCountriesData() expected an error for empty JSON array, but got nil")
		}
		expectedErrorMsg := "apiResponse slice is empty"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getRestCountriesData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 5: Network Error ===
	t.Run("Network Error", func(t *testing.T) {
		// 1. Setup and immediately close server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This handler should never be called during this test case
			t.Errorf("UNEXPECTED: Network error test server recieved a request: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		closedServerURL := server.URL + APIString
		server.Close()

		// 2. Use a standard client that will attempt connection
		testClient := &http.Client{}

		// 3. Call the function under test
		_, actualErr := getRestCountriesData(testClient, closedServerURL, defaultIsoCode)

		// 4. Assertions
		if actualErr == nil {
			t.Fatalf("getRestCountriesData() expected a network error when connecting to %s, but got nil", closedServerURL)
		}
		expectedErrorSubstrings := []string{"connection refused", "connect: connection refused"}
		errorMatched := false
		errStr := actualErr.Error() // Get the error message string

		// 5. Iterate trough expected errors and compare
		for _, sub := range expectedErrorSubstrings {
			if strings.Contains(errStr, sub) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			t.Errorf("getRestCountriesData() error = %q, did not contain expected network error substrings (%v)", actualErr, expectedErrorSubstrings)
		}
	})
}

/*
*	Function with tests for getMetroData()
 */
func TestGetMetroData(t *testing.T) {
	// Define standard inputs used across multiple tests
	defaultLat := 62
	defaultLong := 10
	APIString := "/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation"

	// === Test Case 1: Success Path ===
	t.Run("Success metro data for Norway", func(t *testing.T) {
		// 1. Setting up tests with mock JSON response.
		jsonPath := filepath.Join("testdata", "MetroNorwaySuccess.json")
		mockJSONResponse, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("TEST SETUP FAILED: Could not read testfile %s: %v", jsonPath, err)
		}

		// 2. Set up test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(mockJSONResponse)
		}))
		defer server.Close()

		// 3. Define the expected response
		expectedData := utils.MetroMeanValues{
			MeanPrecipitation: 0.03,
			MeanTemperature:   0.86,
		}

		// 4. Call the function under test
		testClient := server.Client()
		actualData, actualErr := getMetroData(testClient, server.URL+APIString, float64(defaultLat), float64(defaultLong))

		// 5. Assertions
		if actualErr != nil {
			t.Fatalf("getMetroData() returned an unexpected error: %v", actualErr)
		}
		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("getMetroData() returned unexpected data.\nGot:\n%#v\nWant:\n%#v", actualData, expectedData)
		}
	})

	// === Test Case 2: API Returns 404 Status ==
	t.Run("API returns 404", func(t *testing.T) {
		// 1. Setup test server to return 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getMetroData(testClient, server.URL+APIString, float64(defaultLat), float64(defaultLong))

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getMetroData() expected an error for 404 status, but got nil")
		}
		expectedErrorMsg := fmt.Sprintf("API returned non-200 status code: %d", http.StatusNotFound)
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getMetroData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 3: Malformed JSON Response ===
	t.Run("Malformed JSON Response", func(t *testing.T) {
		// 1. Setup test server to return bad JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application-json")
			// Invalid JSON
			_, _ = w.Write([]byte(`{"hourly": {"precipitation": [0.1, 0.2,], "temperature_2m": [-1.0, 0.5}}`))
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getMetroData(testClient, server.URL+APIString, float64(defaultLat), float64(defaultLong))

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getMetroData() expected an error for malformed JSON, but got nil")
		}
		expectedErrorMsg := "error when decoding json"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getMetroData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 4: Network Error ===
	t.Run("Network Error", func(t *testing.T) {
		// 1. Setup and immediately close server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This handler should never be called during this test case
			t.Errorf("UNEXPECTED: Network error test server recieved a request: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		closedServerURL := server.URL + APIString
		server.Close()

		// 2. Use a standard client that will attempt connection
		testClient := &http.Client{}

		// 3. Call the function under test
		_, actualErr := getMetroData(testClient, closedServerURL, float64(defaultLat), float64(defaultLong))

		// 4. Assertions
		if actualErr == nil {
			t.Fatalf("getMetroData() expected a network error when connecting to %s, but got nil", closedServerURL)
		}
		expectedErrorSubstrings := []string{"connection refused", "connect: connection refused"}
		errorMatched := false
		errStr := actualErr.Error() // Get the error message string

		// 5. Iterate trough expected errors and compare
		for _, sub := range expectedErrorSubstrings {
			if strings.Contains(errStr, sub) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			t.Errorf("getMetroData() error = %q, did not contain expected network error substrings (%v)", actualErr, expectedErrorSubstrings)
		}
	})
}

/*
*	Function with tests for getCurrencyData()
 */
func TestGetCurrencyData(t *testing.T) {
	// Define standard inputs used across multiple tests
	APIString := "/currency/"
	defaultInputCurrencies := map[string]interface{}{
		"NOK": map[string]interface{}{
			"name":   "Norwegian krone",
			"symbol": "kr",
		},
	}
	defaultTargetCurrencies := []string{"EUR", "USD", "SEK", "XYZ"}

	// === Test Case 1: Success Path ===
	t.Run("Success currency data for Norway", func(t *testing.T) {
		// 1. Setting up tests with mock JSON response.
		jsonPath := filepath.Join("testdata", "CurrencyNorwaySuccess.json")
		mockJSONResponse, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("TEST SETUP FAILED: Could not read testfile %s: %v", jsonPath, err)
		}

		// 2. Set up test server to return mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(mockJSONResponse)
		}))
		defer server.Close()

		// 3. Define the expected response
		expectedData := map[string]float64{
			"EUR": 0.08817,
			"USD": 0.095174,
			"SEK": 0.954623,
		}

		// 4. Call the function under test
		testClient := server.Client()
		actualData, actualErr := getCurrencyData(testClient, server.URL+APIString, defaultInputCurrencies, defaultTargetCurrencies)

		// 5. Assertions
		if actualErr != nil {
			t.Fatalf("getCurrencyData() returned an unexpected error: %v", actualErr)
		}
		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("getCurrencyData() returned unexpected data.\nGot:\n%#v\nWant:\n%#v", actualData, expectedData)
		}
	})

	// === Test Case 2: API Returns 404 Status ==
	t.Run("API returns 404", func(t *testing.T) {
		// 1. Setup test server to return 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getCurrencyData(testClient, server.URL+APIString, defaultInputCurrencies, defaultTargetCurrencies)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getCurrencyData() expected an error for 404 status, but got nil")
		}
		expectedErrorMsg := fmt.Sprintf("API returned non-200 status code: %d", http.StatusNotFound)
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getCurrencyData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 3: Malformed JSON Response ===
	t.Run("Malformed JSON Response", func(t *testing.T) {
		// 1. Setup test server to return bad JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application-json")
			// Invalid JSON
			_, _ = w.Write([]byte(`{"rates": {"USD": 1.23`))
		}))
		defer server.Close()

		// 2. Call the function under test
		testClient := server.Client()
		_, actualErr := getCurrencyData(testClient, server.URL+APIString, defaultInputCurrencies, defaultTargetCurrencies)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getCurrencyData() expected an error for malformed JSON, but got nil")
		}
		expectedErrorMsg := "error when decoding json"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getCurrencyData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 4: Empty Input Currency Map

	t.Run("Empty Input Currency Map", func(t *testing.T) {
		// 1. Setup specific input for this case
		inputCurrencies := map[string]interface{}{}
		// No server needed for this test as it fails before HTTP request
		// 2. Call the function under test
		testClient := http.DefaultClient
		_, actualErr := getCurrencyData(testClient, "http://example.com", inputCurrencies, defaultTargetCurrencies)

		// 3. Assertions
		if actualErr == nil {
			t.Fatal("getCurrencyData() expected an error for empty currency map, but got nil")
		}
		expectedErrorMsg := "error: no currencies found for country"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("getCurrencyData() error message = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

	// === Test Case 5: Network Error ===
	t.Run("Network Error", func(t *testing.T) {
		// 1. Setup and immediately close server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This handler should never be called during this test case
			t.Errorf("UNEXPECTED: Network error test server recieved a request: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		closedServerURL := server.URL + APIString
		server.Close()

		// 2. Use a standard client that will attempt connection
		testClient := &http.Client{}

		// 3. Call the function under test
		_, actualErr := getCurrencyData(testClient, closedServerURL, defaultInputCurrencies, defaultTargetCurrencies)

		// 4. Assertions
		if actualErr == nil {
			t.Fatalf("getCurrencyData() expected a network error when connecting to %s, but got nil", closedServerURL)
		}
		expectedErrorSubstrings := []string{"connection refused", "connect: connection refused"}
		errorMatched := false
		errStr := actualErr.Error() // Get the error message string

		// 5. Iterate trough expected errors and compare
		for _, sub := range expectedErrorSubstrings {
			if strings.Contains(errStr, sub) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			t.Errorf("getCurrencyData() error = %q, did not contain expected network error substrings (%v)", actualErr, expectedErrorSubstrings)
		}
	})
}

/*
*	Uses test tables to test the function calculateMean(...)
 */
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
