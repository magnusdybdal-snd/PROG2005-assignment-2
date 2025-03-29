package handlers

import (
	"assignment2/utils"
	"context"
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

	log.Printf("Retrieved dashboard for %s (ISO: %s)", config.Country, config.IsoCode)

	return config, nil
}


