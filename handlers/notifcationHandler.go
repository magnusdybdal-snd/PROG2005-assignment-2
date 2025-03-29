package handlers

//Todo: add pathing to make sure things arent routed wrong for the methods

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var ctx context.Context
var client *firestore.Client

func InitFirestore(fc *firestore.Client) {
	client = fc
}

func HandleNotification(w http.ResponseWriter, r *http.Request) {
	ctx = r.Context()
	switch r.Method {
	case http.MethodPost:
		registerNewWebhook(w, r)
	case http.MethodDelete:
		deleteWebhook(w, r)
	case http.MethodGet:
		getWebhooks(w, r)
	default:

	}
}

func registerNewWebhook(w http.ResponseWriter, r *http.Request) {
	//making sure the path is correct
	if r.URL.Path != utils.NOTIFICATION_PATH {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		output := "Expected path is: " + utils.NOTIFICATION_PATH + "<br>" +
			"Supported methods are: " + http.MethodPost + http.MethodDelete + http.MethodGet
		_, err := fmt.Fprint(w, output)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	//struct to be registerd
	var data utils.RegisterWebhook

	//Setupp the decoder so it follows the schema wished upon
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	//Unknown fields are present
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Unknown fields are present", http.StatusBadRequest)
		return
	}

	//make sure that the data is valid
	if data.Url == "" || checkEvent(data) || len(data.Country) != 2 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//send the struct to the databse
	id, _, err := client.Collection(utils.WebhooksCollection).Add(ctx, data)
	if err != nil {
		log.Println("Error when adding webhook: " + err.Error())
		http.Error(w, "Error when adding document", http.StatusBadRequest)
		return
	}

	//struct to send response
	returnId := utils.WebhookId{
		ID: id.ID,
	}

	//preare and send the ID back to user
	w.Header().Set("content-type", "applications/json")
	if err := json.NewEncoder(w).Encode(returnId); err != nil {
		log.Println("Unable to encode the response ID for a webhook registration: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	//get the document ID
	webhookID := r.PathValue("id")

	//Connect to firebase and get the document
	docRef := client.Collection(utils.WebhooksCollection).Doc(webhookID)

	//test if the document can be opend
	_, err := docRef.Get(ctx)
	if err != nil {
		log.Println("Error fetching document: ", err)
		http.Error(w, "Document does not exist", http.StatusBadRequest)
		return
	}

	//delete the document
	_, err = docRef.Delete(ctx)
	if err != nil {
		log.Println("Error deleting document: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	//prepare succsessfull response
	response := map[string]string{"message": "Successfully deleted webhook", "webhookId": webhookID}
	w.Header().Set("Content-Type", "application/json")

	//write status code and send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getWebhooks(w http.ResponseWriter, r *http.Request) {
	//get the path
	webhookID := r.PathValue("id")

	//end response
	var response interface{}

	//check if the path is empty
	if webhookID == "" {
		var webhooks []utils.ReturnWebhook
		iter := client.Collection(utils.WebhooksCollection).Documents(ctx)
		defer iter.Stop()

		//loop over the documents
		for {
			//get the document
			doc, err := iter.Next()

			//tells the for loop when to break out
			if err == iterator.Done {
				break
			}

			//in case of an unforseen error
			if err != nil {
				log.Println("Error fetching document from Firebase:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			//format the response into a struct
			var webhook utils.ReturnWebhook
			if err := doc.DataTo(&webhook); err != nil {
				log.Println("Error mapping data from Firestore into struct:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			webhook.ID = doc.Ref.ID
			webhooks = append(webhooks, webhook)
		}
		response = webhooks
	} else {
		//start the connection and get the document
		docRef := client.Collection(utils.WebhooksCollection).Doc(webhookID)
		doc, err := docRef.Get(ctx)
		if err != nil {
			log.Println("Tried to get a document that does not exist:", err)
			http.Error(w, "Document does not exist", http.StatusBadRequest)
			return
		}

		//get the struct and decode into it
		var webhook utils.ReturnWebhook
		if err := doc.DataTo(&webhook); err != nil {
			log.Println("Error unmarshaling data into struct:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		//set the ID and response
		webhook.ID = docRef.ID
		response = webhook
	}

	// Send the response in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// used to check if the event is correct
func checkEvent(hook utils.RegisterWebhook) bool {
	switch hook.Event {
	case "REGISTER", "CHANGE", "DELETE", "INVOKE":
		return false
	default:
		return true
	}
}
