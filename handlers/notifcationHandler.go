package handlers

//Todo: add pathing to make sure things arent routed wrong for the methods

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func HandleNotification(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		registerNewWebhook(w, r, false)
	case http.MethodDelete:
		deleteWebhook(w, r, false)
	case http.MethodGet:
		getWebhooks(w, r, false)
	case http.MethodPatch:
		patchWebhook(w, r, false)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func patchWebhook(w http.ResponseWriter, r *http.Request, isTest bool) {
	log.Println("Recieved " + r.Method + " request.")
	//values to be used to talk to firestore
	var webhookID string
	var collection string
	if isTest {
		webhookID = strings.TrimPrefix(r.URL.Path, utils.NOTIFICATION_PATH)
		collection = utils.WEBHOOKTESTCOLLECTION
	} else {
		webhookID = r.PathValue("id")
		collection = utils.WebhooksCollection
	}

	//get context and path value
	ctx := r.Context()
	//check if the path is empty
	if webhookID == "" {
		log.Println("Error, webhook id is required")
		http.Error(w, "Error webhook id is required.", http.StatusBadRequest)
		return
	}

	// Test for embedded webhook id
	if webhookID == "" {
		log.Println("Error, webhook id is required")
		http.Error(w, "Error webhook id is required.", http.StatusBadRequest)
		return
	}

	// Reads body of PATCH request
	var content utils.RegisterWebhook
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&content)
	if err != nil {
		log.Println("Error decoding JSON payload: ", err)
		http.Error(w, "Error in request body", http.StatusBadRequest)
		return
	}

	// if all fields are empty, return 400 error
	if content.Url == "" && content.Country == "" && content.Event == "" {
		http.Error(w, "Invalid empty json request", http.StatusBadRequest)
		log.Println("Tried to update empty json request")
		return
	}

	var update []firestore.Update

	// Adds the wanted fields to the update list
	if content.Url != "" {
		update = append(update, firestore.Update{Path: "url", Value: content.Url})
	}
	if content.Country == "" || len(content.Country) == 2 {
		update = append(update, firestore.Update{Path: "country", Value: strings.ToUpper(content.Country)})
	}
	if !checkEvent(content) {
		update = append(update, firestore.Update{Path: "event", Value: content.Event})
	}

	// Update the document in Firestore
	docRef := utils.FirestoreClient.Collection(collection).Doc(webhookID)
	_, err = docRef.Update(ctx, update)
	if err != nil {
		log.Println("Error updating document: ", err)
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func registerNewWebhook(w http.ResponseWriter, r *http.Request, isTest bool) {
	if r.URL.Path != utils.NOTIFICATION_PATH {
		http.Error(w, "This method does not offer any functionality outside of: "+utils.NOTIFICATION_PATH, http.StatusBadRequest)
		return
	}
	log.Println("Recieved " + r.Method + " request.")

	//if it is a test, you need to set the client manually
	//variables will have to be set differently too
	var collection string
	if isTest {
		collection = utils.WEBHOOKTESTCOLLECTION
	} else {
		collection = utils.WebhooksCollection
	}
	ctx := r.Context()

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

	data.Country = strings.ToUpper(data.Country)

	//make sure that the data is valid
	if data.Url != "" && len(data.Url) == 2 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if checkEvent(data) {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//send the struct to the databse
	id, _, err := utils.FirestoreClient.Collection(collection).Add(ctx, data)
	if err != nil {
		log.Println("Error when adding webhook: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//struct to send response
	returnId := utils.WebhookId{
		ID: id.ID,
	}

	responseJson, err := json.Marshal(returnId)
	if err != nil {
		log.Println("Unable to encode the response ID for a webhook registration: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//preare and send the ID back to user
	w.Header().Set("content-type", "applications/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}

func deleteWebhook(w http.ResponseWriter, r *http.Request, isTest bool) {
	log.Println("Recieved " + r.Method + " request.")

	//variables will be set according to the test or not
	var collection string
	var webhookID string

	//if it is a test, you need to set the client manually
	//variables will have to be set differently too
	if isTest {
		webhookID = strings.TrimPrefix(r.URL.Path, utils.NOTIFICATION_PATH)
		collection = utils.WEBHOOKTESTCOLLECTION
	} else {
		collection = utils.WebhooksCollection
		webhookID = r.PathValue("id")
	}

	ctx := r.Context()

	//Connect to firebase and get the document
	docRef := utils.FirestoreClient.Collection(collection).Doc(webhookID)

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
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}

func getWebhooks(w http.ResponseWriter, r *http.Request, isTest bool) {
	log.Println("Recieved " + r.Method + " request.")

	//variables will be set according to the test or not
	var collection string
	var webhookID string

	//if it is a test, you need to set the client manually
	//variables will have to be set differently too
	if isTest {
		//not mapped trough main, does not know what {id} is
		webhookID = strings.TrimPrefix(r.URL.Path, utils.NOTIFICATION_PATH)
		collection = utils.WEBHOOKTESTCOLLECTION
	} else {
		webhookID = r.PathValue("id")
		collection = utils.WebhooksCollection
	}
	ctx := r.Context()
	//get the path

	//end response
	var response interface{}

	//check if the path is empty
	if webhookID == "" {
		var webhooks []utils.ReturnWebhook
		iter := utils.FirestoreClient.Collection(collection).Documents(ctx)
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
			//no point in sending a response to the webhooks if you cannot get the data
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
		docRef := utils.FirestoreClient.Collection(collection).Doc(webhookID)
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

func retriveWebhooks(ctx context.Context) ([]utils.ReturnWebhook, error) {

	var webhooks []utils.ReturnWebhook
	iter := utils.FirestoreClient.Collection(utils.WebhooksCollection).Documents(ctx)
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
			return nil, err
		}
		//format the response into a struct
		var webhook utils.ReturnWebhook
		if err := doc.DataTo(&webhook); err != nil {
			log.Println("Error mapping data from Firestore into struct:", err)
			return nil, err
		}
		webhook.ID = doc.Ref.ID
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}

func invokeWebhook(event string, countryIso2 string, ctx context.Context) {
	webhooks, err := retriveWebhooks(ctx)
	log.Println("invokeWebhook invoked, event: ", event)
	if err != nil {
		log.Println("Error in retrieving webhooks ", err)
		return
	}
	for _, v := range webhooks {
		if (v.Country == "" || v.Country == countryIso2) && v.Event == event {
			log.Println("Activated webhook: ", v.ID, " event: ", v.Event)
			go utils.CallUrl(v.Url, utils.INVOKE, v)
		}
	}
}

// used to check if the event is correct
func checkEvent(hook utils.RegisterWebhook) bool {
	switch hook.Event {
	case utils.REGISTER, utils.CHANGE, utils.DELETE, utils.INVOKE, utils.NOTREACHABLE:
		return false
	default:
		return true
	}
}
