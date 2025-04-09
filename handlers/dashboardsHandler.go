package handlers

/*
*	TODO: 	Add timeouts for API calls
	TODO:	Consider adding context to http request for proper timeout management
	TODO: 	Currently calls to dashboards/ and dashboards/ID/something goes to default handler
*/

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

/*
*	These variables are used to swap the actual functions to mocked ones during testing
 */
var (
	getConficFunc             = utils.GetFirestoreDocument[utils.DashboardConfig]
	getCountriesFunc          = getRestCountriesData
	getMetroFunc              = getMetroData
	getCurrencyFunc           = getCurrencyData
	tryCacheRestCountriesFunc = tryCacheRestCountries
	tryCacheMetroFunc         = tryCacheMetro
	tryCacheCurrencyFunc      = tryCacheCurrency
)

func HandleGetDashboard(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Test for embedded dashboard id
	dashboardId := r.PathValue("id")
	if dashboardId == "" {
		log.Println("Error, dashboard id is required")
		http.Error(w, "Dashboard id is required.", http.StatusBadRequest)
		return
	}

	// Retrieves the dashboard configuration from firestore database
	dashboardConfig, err := getConficFunc(ctx, dashboardId, utils.DASHBOARD_COLLECTION)
	if err != nil {
		log.Printf("Error retrieving dashboard config from database: %v", err)
		http.Error(w, "Could not find dashboard "+dashboardId, http.StatusInternalServerError)
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
		restCountriesData, err = tryCacheRestCountries(ctx, http.DefaultClient, utils.RESTCountriesAPI, dashboardConfig.IsoCode)
		if err != nil {
			log.Printf("Error getting RestCountries data: %v", err)
			http.Error(w, "Error getting country information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Metro API if needed
	if needMetroAPI {
		metroData, err = tryCacheMetro(ctx, http.DefaultClient, utils.MetroAPI, float64(restCountriesData.Coordinates[0]), float64(restCountriesData.Coordinates[1]))
		if err != nil {
			log.Printf("Error getting MetroAPI data: %v", err)
			http.Error(w, "Error getting weather information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Currency API if needed
	if needCurrencyAPI {
		currencyData, err = tryCacheCurrency(ctx, http.DefaultClient, utils.CurrencyAPI, restCountriesData.Currencies, dashboardConfig.Features.TargetCurrencies)
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
	invokeWebhook(utils.INVOKE, response.IsoCode, ctx)
}

func getRestCountriesData(ctx context.Context, client *http.Client, baseURL string, IsoCode string) (utils.RestCountriesResponse, error) {

	// Url to invoke
	url := baseURL + IsoCode

	// Uses get with default http client
	resp, err := client.Get(url)
	if err != nil {
		return utils.RestCountriesResponse{}, fmt.Errorf("error fetching country info form REST Countries: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", ctx)
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
func getMetroData(ctx context.Context, client *http.Client, baseURL string, lat float64, long float64) (utils.MetroMeanValues, error) {

	// Url to invoke
	url := fmt.Sprintf(baseURL, lat, long)

	// Uses http.Get with standard client and does the request
	resp, err := client.Get(url)
	if err != nil {
		return utils.MetroMeanValues{}, fmt.Errorf("error fetching weather data from Metro API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", ctx)
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

func getCurrencyData(ctx context.Context, client *http.Client, baseURL string, currencies map[string]interface{}, targetCurrencies []string) (map[string]float64, error) {
	// Extracts the FIRST currency if there are more than one
	var currencyISO string
	for iso := range currencies {
		currencyISO = iso
		break
	}
	// Checks that a currency is found for the country
	if currencyISO == "" {
		return nil, fmt.Errorf("error: no currencies found for country")
	}
	// url to invoke
	url := baseURL + currencyISO

	// Uses http.Get with standard client and does the request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching currency data from Currency API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			invokeWebhook(utils.NOTREACHABLE, "", ctx)
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

func tryCacheRestCountries(ctx context.Context, client *http.Client, baseURL string, IsoCode string) (utils.RestCountriesResponse, error) {
	cacheKey := fmt.Sprintf("restcountries_%s", IsoCode)
	docRef := utils.FirestoreClient.Collection(utils.CACHE_COLLECTION).Doc(cacheKey)

	doc, err := docRef.Get(ctx)
	if err == nil { // cache hit
		var cachedData utils.CachedRestCountries
		if err2 := doc.DataTo(&cachedData); err == nil {
			return cachedData.Data, nil // Sends the cahced data back
		} else {
			log.Printf("Error unmarshalling cached RestCountries data: %v", err2)
		}
	} else {
		log.Printf("Cache miss for RestCountries: %s", IsoCode)
	}

	// Cache miss or error: Get new data
	newData, newErr := getCountriesFunc(ctx, client, baseURL, IsoCode)
	if newErr != nil {
		return utils.RestCountriesResponse{}, newErr
	}

	// Store the new data in cache
	cacheEntry := utils.CachedRestCountries{
		Data:      newData,
		Timestamp: time.Now(),
	}
	_, addErr := docRef.Set(ctx, cacheEntry)
	if addErr != nil {
		log.Printf("Failed to save RestCountries to cache: %v", addErr)
	}

	return newData, nil
}

func tryCacheMetro(ctx context.Context, client *http.Client, baseURL string, lat float64, long float64) (utils.MetroMeanValues, error) {
	cacheKey := fmt.Sprintf("metro_%f_%f", lat, long)
	docRef := utils.FirestoreClient.Collection(utils.CACHE_COLLECTION).Doc(cacheKey)

	doc, err := docRef.Get(ctx)
	if err == nil { // cache hit
		var cachedData utils.CachedMetro
		if err2 := doc.DataTo(&cachedData); err == nil {
			return cachedData.Data, nil // Sends the cahced data back
		} else {
			log.Printf("Error unmarshalling cached Metro data: %v", err2)
		}
	} else {
		log.Printf("Cache miss for Metro: long:%f lat:%f", long, lat)
	}

	// Cache miss or error: Get new data
	newData, newErr := getMetroFunc(ctx, client, baseURL, lat, long)
	if newErr != nil {
		return utils.MetroMeanValues{}, newErr
	}

	// Store the new data in cache
	cacheEntry := utils.CachedMetro{
		Data:      newData,
		Timestamp: time.Now(),
	}
	_, addErr := docRef.Set(ctx, cacheEntry)
	if addErr != nil {
		log.Printf("Failed to save Metro to cache: %v", addErr)
	}

	return newData, nil
}

func tryCacheCurrency(ctx context.Context, client *http.Client, baseURL string, currencies map[string]interface{}, targetCurrencies []string) (map[string]float64, error) {

	var currencyISO string
	for iso := range currencies {
		currencyISO = iso
		break
	}
	cacheKey := fmt.Sprintf("metro_%s", currencyISO)
	docRef := utils.FirestoreClient.Collection(utils.CACHE_COLLECTION).Doc(cacheKey)

	doc, err := docRef.Get(ctx)
	if err == nil { // cache hit
		var cachedData utils.CachedCurrency
		if err2 := doc.DataTo(&cachedData); err == nil {
			return cachedData.Data, nil // Sends the cahced data back
		} else {
			log.Printf("Error unmarshalling cached Currency data: %v", err2)
		}
	} else {
		log.Printf("Cache miss for Currency: %s", currencyISO)
	}

	// Cache miss or error: Get new data
	newData, newErr := getCurrencyFunc(ctx, client, baseURL, currencies, targetCurrencies)
	if newErr != nil {
		return nil, newErr
	}

	// Store the new data in cache
	cacheEntry := utils.CachedCurrency{
		Data:      newData,
		Timestamp: time.Now(),
	}
	_, addErr := docRef.Set(ctx, cacheEntry)
	if addErr != nil {
		log.Printf("Failed to save Currency to cache: %v", addErr)
	}

	return newData, nil
}

/*
*	Function that returns the mean value, calculated from a list of values
 */
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
