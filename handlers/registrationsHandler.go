package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// TODO legge inn fprintf.log for error logging i stedet for log.Println

var ctx context.Context
var client *firestore.Client

// InitFirestore allows main to pass the Firestore client to handlers
func InitFirestore(fc *firestore.Client, context context.Context) {
	client = fc
	ctx = context
}

func HandleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Not yet implemented " + r.Method)
	case http.MethodPost:
		registerDashConfig(w, r) // Ensures that registration of dashboard config handles POST requests.
	default:
		log.Println("Unsupported method " + r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func registerDashConfig(w http.ResponseWriter, r *http.Request) {
	// Initial client error handling
	if client == nil {
		log.Println("ERROR: Firestore client is nil - not initialized")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	//log.Println("Starting registerDashConfig handler") // debug log
	//log.Println("Content-Type:", r.Header.Get("Content-Type")) // debug log

	content, err := io.ReadAll(r.Body) // TODO read payload and check that it is up to spec.
	if err != nil {
		log.Println("Reading payload from body failed:", err)
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}

	//log.Println("Request body:", string(content)) // debug log

	if len(string(content)) == 0 {
		//log.Println("Content appears to be empty.") // debug log
		http.Error(w, "Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	} else {
		//log.Println("Unmarshalling JSON") // debug log
		s := utils.DashboardConfig{}
		err := json.Unmarshal(content, &s)
		if err != nil {
			log.Println("Error unmarshalling payload:", err)
			http.Error(w, "Error unmarshalling payload: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//log.Println("Unmarshalled successfully, adding to Firestore") // debug log

		s.LastRetrieval = time.Now() // update timestamp

		id, _, err2 := client.Collection(utils.DASHBOARD_COLLECTION).Add(ctx, s)
		if err2 != nil {
			log.Println("Error when adding document:", err2)
			http.Error(w, "Error when adding document: "+err2.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Document added successfully, creating response")
			response := struct {
				ID            string    `json:"id"`
				LastRetrieval time.Time `json:"lastChange"`
			}{
				ID:            id.ID,
				LastRetrieval: s.LastRetrieval,
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				log.Println("Error creating JSON response:", err)
				http.Error(w, "Error creating response", http.StatusInternalServerError)
				return
			}

			//log.Println("Setting headers and writing response") // debug log
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			n, err := w.Write(responseJSON) // This sould be fine as long as the json is checked properly ln53.
			if err != nil {
				log.Println("Error writing response:", err)
			} else {
				log.Println("Wrote", n, "bytes successfully")
			}
			//log.Println("Handler completed successfully") // debug log
			return
		}
	}
}
