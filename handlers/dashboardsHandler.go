package handlers

/*
*	TODO: 	Add timeouts for API calls
	TODO:	Consider adding context to http request for proper timeout management
*/

import (
	"assignment2/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

func HandleGetDashboard(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// Test for embedded dashboard id
	dashboardId := r.PathValue("id")
	if dashboardId == "" {
		log.Println("Error, dashboard id is required")
		http.Error(w, "error dashboard id is required.", http.StatusBadRequest)
		return
	}

	// Retrieves the dashboard configuration from firestore database
	dashboardConfig, err := utils.GetDashboardConfig[utils.DashboardConfig](ctx, dashboardId, utils.DASHBOARD_COLLECTION)
	if err != nil {
		log.Printf("Error retrieving dashboard config from database: %v", err)
		http.Error(w, "error retrieving dashboard", http.StatusInternalServerError)
		return
	}

	// Find out which API calls we need to do.
	// Both Metro and Currency need information from RestCountries to be invoked
	needMetroAPI := dashboardConfig.Features.Temperature || dashboardConfig.Features.Precipitation

	needCurrencyAPI := len(dashboardConfig.Features.TargetCurrencies) > 0

	needRestCountriesAPI := dashboardConfig.Features.Capital || dashboardConfig.Features.Coordinates ||
		dashboardConfig.Features.Area || dashboardConfig.Features.Population ||
		needMetroAPI || needCurrencyAPI

	var restCountriesData utils.RestCountriesResponse
	var metroData utils.MetroMeanValues
	var currencyData map[string]float64

	// Gets the data from REST Countries if needed
	if needRestCountriesAPI {
		restCountriesData, err = getRestCountriesData(dashboardConfig.IsoCode, r)
		if err != nil {
			log.Printf("Error getting RestCountries data: %v", err)
			http.Error(w, "Error getting country information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Metro API if needed
	if needMetroAPI {
		metroData, err = getMetroData(float64(restCountriesData.Coordinates[0]), float64(restCountriesData.Coordinates[1]), r)
		if err != nil {
			log.Printf("Error getting MetroAPI data: %v", err)
			http.Error(w, "Error getting weather information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Currency API if needed
	if needCurrencyAPI {
		currencyData, err = getCurrencyData(restCountriesData.Currencies, dashboardConfig.Features.TargetCurrencies, r)
		if err != nil {
			log.Printf("Error getting Currency API data: %v", err)
			http.Error(w, "Error getting currency information", http.StatusInternalServerError)
			return
		}
	}

	var response utils.DashboardResponse
	response.Country = dashboardConfig.Country
	response.IsoCode = dashboardConfig.IsoCode
	response.LastRetrieval = time.Now()

	if dashboardConfig.Features.Capital {
		response.Features.Capital = restCountriesData.Capital[0]
	}

	if dashboardConfig.Features.Area {
		response.Features.Area = restCountriesData.Area
	}

	if dashboardConfig.Features.Population {
		response.Features.Population = restCountriesData.Population
	}

	if dashboardConfig.Features.Coordinates {
		response.Features.Coordinates = map[string]float64{
			"latitude":  float64(restCountriesData.Coordinates[0]),
			"longitude": float64(restCountriesData.Coordinates[1]),
		}
	}

	if dashboardConfig.Features.Temperature {
		response.Features.Temperature = metroData.MeanTemperature
	}

	if dashboardConfig.Features.Precipitation {
		response.Features.Precipitation = metroData.MeanPrecipitation
	}

	if needCurrencyAPI {
		response.Features.TargetCurrencies = currencyData
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	invokeWebhook(utils.INVOKE, response.IsoCode, r)
}

func getRestCountriesData(IsoCode string, r *http.Request) (utils.RestCountriesResponse, error) {

	// Url to invoke
	url := utils.RESTCountriesAPI + IsoCode

	// Uses http.Get to setup standard client and do the request
	resp, err := http.Get(url)
	if err != nil {
		return utils.RestCountriesResponse{}, fmt.Errorf("error fetching country info form REST Countries: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", r)
		}
		return utils.RestCountriesResponse{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Decodes json response into struct. REST Countries always returns an array of countries, even tho we only ask for one
	var apiResponse []utils.RestCountriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return utils.RestCountriesResponse{}, fmt.Errorf("error when decoding json: %v", err)
	}

	// Checks that the apiResponse slice is not empty
	if len(apiResponse) == 0 {
		return utils.RestCountriesResponse{}, fmt.Errorf("apiResponse slice is empty")
	}

	return apiResponse[0], nil
}

/*
*	This function invokes the Metro API with the parameter latitude and logitude, and returns temperature and precipiation hourly
*	for a 7 day forecast as a struct with two lists. TODO: calculate mean value and return the mean values as a list??
 */
func getMetroData(lat float64, long float64, r *http.Request) (utils.MetroMeanValues, error) {

	// Url to invoke
	url := fmt.Sprintf(utils.MetroAPI, lat, long)

	// Uses http.Get with standard client and does the request
	resp, err := http.Get(url)
	if err != nil {
		return utils.MetroMeanValues{}, fmt.Errorf("error fetching weather data from Metro API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", r)
		}
		return utils.MetroMeanValues{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Decodes json response into struct.
	var apiResponse utils.MetroResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return utils.MetroMeanValues{}, fmt.Errorf("error when decoding json: %v", err)
	}

	meanPrecip := calculateMean(apiResponse.Hourly.Precipitation)
	meanTemp := calculateMean(apiResponse.Hourly.Temperature2M)

	meanResponse := utils.MetroMeanValues{
		MeanPrecipitation: meanPrecip,
		MeanTemperature:   meanTemp,
	}

	return meanResponse, nil
}

func getCurrencyData(currencies map[string]interface{}, targetCurrencies []string, r *http.Request) (map[string]float64, error) {
	// Extracts the FIRST currency if there are more than one
	var currencyISO string
	for iso := range currencies {
		currencyISO = iso
		break
	}
	// Checks that a currency is found in the response
	if currencyISO == "" {
		return nil, fmt.Errorf("error: no currencies found in response")
	}
	// url to invoke
	url := utils.CurrencyAPI + currencyISO

	// Uses http.Get with standard client and does the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching currency data from Currency API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", r)
		}
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}
	// Decodes json response into struct.
	var apiResponse utils.CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error when decoding json: %v", err)
	}

	filteredRates := make(map[string]float64)
	// Adds the wanted currencies from the API response to the map presented in the response
	for _, currency := range targetCurrencies {
		if rate, exists := apiResponse.Rates[currency]; exists {
			filteredRates[currency] = rate
		}
	}

	return filteredRates, nil
}

func calculateMean(val []float64) float64 {

	if len(val) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range val {
		sum += v
	}
	mean := sum / float64(len(val))
	return math.Round(mean*100) / 100
}
