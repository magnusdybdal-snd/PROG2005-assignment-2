package handlers

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
)

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
	
}


