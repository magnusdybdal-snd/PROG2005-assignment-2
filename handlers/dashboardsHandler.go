package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

const collection = "dashboards"

type DashboardHandler struct {
	fsClient *firestore.Client
}

func NewDashboardHandler (fsClient *firestore.Client) *DashboardHandler {
	return &DashboardHandler{fsClient: fsClient}
}

func (h *DashboardHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.handleGetDashboard(ctx, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DashboardHandler) handleGetDashboard (ctx context.Context, w http.ResponseWriter, r *http.Request) {

	// Test for embedded dashboard id
	dashboardId := r.PathValue("id")
	if dashboardId == "" {
		log.Println("Error, dashboard id is required")
		http.Error(w, "error dashboard is is required.", http.StatusBadRequest)
		return
	}

	// Retrieves the dashboard configuration from firestore database
	dashboardConfig, err := h.getDashboardConfig(ctx, dashboardId)
	if err != nil {
		log.Printf("Error retrieving dashboard config from database: %v", err)
		http.Error(w, "error retrieving dashboard", http.StatusInternalServerError)
		return
	}

	// Find out which API calls we need to do.
	// Both Metro and Currency need information from RestCountries to be invoked
	needMetroAPI :=			dashboardConfig.Features.Temperature || dashboardConfig.Features.Precipiation

	needCurrencyAPI :=		len(dashboardConfig.Features.TargetCurrencies) > 0
	
	needRestCountriesAPI := dashboardConfig.Features.Capital     || dashboardConfig.Features.Coordinates ||
						    dashboardConfig.Features.Area        || dashboardConfig.Features.Population ||
						    needMetroAPI						 || needCurrencyAPI


	var restCountriesData utils.RestCountriesResponse
	var metroData utils.MetroResponse
	var currencyData utils.CurrencyResponse

	// REMOVE WHEN DONE
	fmt.Printf("DEBUG: %v %v", metroData, currencyData)

	// Gets the data from REST Countries if needed
	if needRestCountriesAPI {
		restCountriesData, err = h.getRestCountriesData(dashboardConfig.IsoCode)
		if err != nil {
			log.Printf("Error getting RestCountries data: %v", err)
			http.Error(w, "Error getting country information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Metro API if needed
	if needMetroAPI {
		_, err := h.getMetroData(restCountriesData.Coordinates[0], restCountriesData.Coordinates[1])
		if err != nil {
			log.Printf("Error getting MetroAPI data: %v", err)
			http.Error(w, "Error getting weather information", http.StatusInternalServerError)
			return
		}
	}

	// Gets the data from Currency API if needed
	if needCurrencyAPI {
		_, err := h.getCurrencyData(restCountriesData.Currencies)
		if err != nil {
			log.Printf("Error getting Currency API data: %v", err)
			http.Error(w, "Error getting currency information", http.StatusInternalServerError)
			return
		}
	}
	
	

}

func (h *DashboardHandler) getDashboardConfig (ctx context.Context, dashboardId string) (utils.DashboardConfig, error) {

	// Gets reference to firestore document
	docRef := h.fsClient.Collection(collection).Doc(dashboardId)

	// Fetches the data from the document
	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Println("Failed to get dashboard document %s, %v" + dashboardId, err)
		return utils.DashboardConfig{}, err
	}

	// Unmarshals the data in the document to struct
	var config utils.DashboardConfig
	err2 := doc.DataTo(&config)
	if err2 != nil {
		log.Printf("Failed to unmarshal dashboard config: %s: %v", dashboardId, err)
		return utils.DashboardConfig{}, err
	}

	return config, nil
}

func (h *DashboardHandler) getRestCountriesData (IsoCode string) (utils.RestCountriesResponse, error) {

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
		return utils.RestCountriesResponse{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Decodes json response into struct. REST Countries always returns an array of countries, even tho we only ask for one
	var apiResponse []utils.RestCountriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return utils.RestCountriesResponse{}, fmt.Errorf("error when decoding json: %v", err)
	}

	// Checks that the apiResponse slice is not empty
	if len(apiResponse) == 0 {
		return utils.RestCountriesResponse{}, fmt.Errorf("apiResponse slice is empty: %v", err)
	}

	return apiResponse[0], nil
}

/*
*	This function invokes the Metro API with the parameter latitude and logitude, and returns temperature and precipiation hourly
*	for a 7 day forecast as a struct with two lists. TODO: calculate mean value and return the mean values as a list?? 
*/
func (h* DashboardHandler) getMetroData (lat int, long int) (utils.MetroResponse, error) {

	// Url to invoke
	url := fmt.Sprintf(utils.MetroAPI, lat, long)

	// Uses http.Get with standard client and does the request
	resp, err := http.Get(url)
	if err != nil {
		return utils.MetroResponse{}, fmt.Errorf("error fetching weather data from Metro API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return utils.MetroResponse{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Decodes json response into struct. REST Countries always returns an array of countries, even tho we only ask for one
	var apiResponse utils.MetroResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return utils.MetroResponse{}, fmt.Errorf("error when decoding json: %v", err)
	}

	return apiResponse, nil
}

func (h* DashboardHandler) getCurrencyData (currencies map[string]interface{}) (utils.CurrencyResponse, error) {
	// Extracts the FIRST currency if there are more than one 
	var currencyISO string
	for iso := range currencies {
		currencyISO = iso
		break
	}
	// Checks that a currency is found in the response
	if currencyISO == "" {
		return utils.CurrencyResponse{}, fmt.Errorf("error: no currencies found in response")
	}
	// url to invoke
	url := utils.CurrencyAPI + currencyISO

	// Uses http.Get with standard client and does the request
	resp, err := http.Get(url)
	if err != nil {
		return utils.CurrencyResponse{}, fmt.Errorf("error fetching currency data from Currency API: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return utils.CurrencyResponse{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}
	// Decodes json response into struct.
	var apiResponse utils.CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return utils.CurrencyResponse{}, fmt.Errorf("error when decoding json: %v", err)
	}

	return apiResponse, nil
}