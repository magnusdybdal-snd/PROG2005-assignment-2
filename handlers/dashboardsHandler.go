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

	// Find out which API calls we need to do
	needRestCountriesAPI := dashboardConfig.Features.Capital     || dashboardConfig.Features.Coordinates ||
						    dashboardConfig.Features.Area        || dashboardConfig.Features.Population ||
						    dashboardConfig.Features.Temperature || dashboardConfig.Features.Precipiation
	
	needMetroAPI :=			dashboardConfig.Features.Temperature || dashboardConfig.Features.Precipiation

	needCurrencyAPI :=		len(dashboardConfig.Features.TargetCurrencies) > 0

	fmt.Printf("DEBUG: needMetroAPI=%v, needCurrencyAPI=%v\n", needMetroAPI, needCurrencyAPI)

	// Gets the data from REST Countries if needed
	if needRestCountriesAPI {
		_, err := h.getRestCountriesData(ctx, dashboardConfig.IsoCode)
		if err != nil {
			log.Printf("Error getting RestCountries data: %v", err)
			http.Error(w, "Error getting country information", http.StatusInternalServerError)
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

func (h *DashboardHandler) getRestCountriesData (ctx context.Context, IsoCode string) (utils.RestCountriesResponse, error) {

	url := utils.RESTCountriesAPI + IsoCode

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch country info from REST Countries: %v", err)
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