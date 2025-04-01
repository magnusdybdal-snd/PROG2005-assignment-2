package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
)

/*
* Handle different types of requests.
 */
func HandleMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		displayDocument(w, r, ctx)
	case http.MethodPost:
		registerDashConfig(w, r, ctx) // Ensures that registration of dashboard config handles POST requests.
	case http.MethodDelete:
		deleteDocument(w, r, ctx)
	default:
		log.Println("Unsupported method " + r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

/*
* Deletes documents using {"id"} in DELETE request.
 */
func deleteDocument(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	log.Println("Received " + r.Method + " request.")

	// Extract id from URL.
	messageId := r.PathValue("id")

	// Retrieve specific message based on id (Firestore-generated hash)
	res := utils.FirestoreClient.Collection(utils.DASHBOARD_COLLECTION).Doc(messageId)

	// Checks if the document exists in database.
	_, err := res.Get(ctx)
	if err != nil {
		log.Println("Document ID does not exist. Id: " + messageId)
		http.Error(w, http.StatusText(http.StatusBadRequest)+": Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve reference to document.
	_, err2 := res.Delete(ctx)
	if err2 != nil {
		log.Println("Delete request for document " + messageId + " failed.")
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
* Registers Dashboard configuratuins in firestore.
 */
func registerDashConfig(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// Initial client error handling
	if utils.FirestoreClient == nil {
		log.Println("ERROR: Firestore client is nil - not initialized")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	log.Println("Starting registerDashConfig handler")
	log.Println("Content-Type:", r.Header.Get("Content-Type"))

	content, err := io.ReadAll(r.Body) // TODO read payload and check that it is up to spec.
	if err != nil {
		log.Println("Reading payload from body failed:", err)
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}

	log.Println("Request body:", string(content))

	if len(string(content)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	} else {
		log.Println("Unmarshalling JSON")
		s := utils.DashboardConfig{}
		err := json.Unmarshal(content, &s)
		if err != nil {
			log.Println("Error unmarshalling payload:", err)
			http.Error(w, "Error unmarshalling payload: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Unmarshalled successfully, adding to Firestore")

		s.LastRetrieval = time.Now() // update timestamp

		id, _, err2 := utils.FirestoreClient.Collection(utils.DASHBOARD_COLLECTION).Add(ctx, s)
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

			log.Println("Setting headers and writing response")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			n, err := w.Write(responseJSON) // This sould be fine as long as the json is checked properly ln53.
			if err != nil {
				log.Println("Error writing response:", err)
			} else {
				log.Println("Wrote", n, "bytes successfully")
			}
			log.Println("Handler completed successfully")
			return
		}
	}
}

/*
* Reads a string from the body in plain-text and sends it to Firestore to be registered as a document.
 */
func displayDocument(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	var response interface{}
	log.Println("Received " + r.Method + " request.")

	// Test for embedded message ID
	messageId := r.PathValue("id")

	// ID id provided in URL
	if messageId != "" {

		// Retrieve specific message based on id (Firestore-generated hash)
		res := utils.FirestoreClient.Collection(utils.DASHBOARD_COLLECTION).Doc(messageId)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Document ID does not exist. Id: " + messageId)
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Invalid ID", http.StatusBadRequest)
			return
		}

		// Creating instance of stuct to be sent.
		var documentResponse utils.RegistrationGetResponse
		documentResponse.Id = messageId // Add document ID to struct.
		if err := doc.DataTo(&documentResponse); err != nil {
			log.Println("Failed to construct struct response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		response = documentResponse

		// No ID in URL: Send all documents.
	} else {
		iter := utils.FirestoreClient.Collection(utils.DASHBOARD_COLLECTION).Documents(ctx)

		// Array used to store multiple documents.
		var registrations []utils.RegistrationGetResponse

		// Loop over documents.
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Printf("failed to iterate: %v", err)
				return
			}

			// Populate struct with values.
			var documentResponse utils.RegistrationGetResponse
			documentResponse.Id = doc.Ref.ID // Add document ID to main.
			if err := doc.DataTo(&documentResponse); err != nil {
				log.Println("Failed to construct struct response")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			// Append each document into array.
			registrations = append(registrations, documentResponse)
		}

		response = registrations
	}
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Encoding of struct failed. ")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
