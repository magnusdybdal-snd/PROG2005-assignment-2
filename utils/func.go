package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CallUrl(url string, event string, content ReturnWebhook) {
	currentTime := time.Now()
	invoke := SendNotification{
		ID:      content.ID,
		Country: content.Country,
		Event:   content.Event,
		Time:    currentTime.Format("20060102 15:04"),
	}
	jsonData, err := json.Marshal(invoke)
	if err != nil {
		log.Println("Error in encoding JSON for webhook call ", err)

	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error during request creation: ", err)
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: ", err)
		return
	}
	log.Println("Webhook " + url + " invoked, recieved status code " + strconv.Itoa(res.StatusCode))
}

/*
*	Function that returns any firestore document, given its collection name and ID.
*
*	param ctx - context used by the handler calling the function
*	param docId - ID of the firestore document
*	param collection - Name of the firestroe collection
*
*	return T - The firestore document that is retrieved
*	return error - error if document cannot be retrieved
 */
func GetDashboardConfig[T any](ctx context.Context, docId string, collection string) (T, error) {

	// Initiate zero value of type T
	var result T

	// Gets reference to firestore document
	docRef := FirestoreClient.Collection(collection).Doc(docId)

	// Fetches the data from the document
	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Println("Failed to get dashboard document %s: %v"+docId, err)
		return result, err
	}

	// Unmarshals the data in the document to struct
	err2 := doc.DataTo(&result)
	if err2 != nil {
		log.Printf("Failed to unmarshal dashboard config: %s: %v", docId, err)
		return result, err
	}

	return result, nil
}

/*
*  Function validates incoming request with struct.
*
*	param valid_schema - Struct used to validate request body.
*	param w 				 - Response writer.
*	param r  			 - Request sent by user.
*
*	return r.Body - sucsessfully validated content
*	return err    - error if post request does not match structs field.
 */
func ValidatePostRequest(valid_schema interface{}, w http.ResponseWriter, r *http.Request) ([]byte, error) {

	// Read request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return nil, err
	}

	// Setup the decoder so it follows the schema wished upon
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()

	// Unknown fields are present - return 400 Error
	if err := decoder.Decode(valid_schema); err != nil {
		http.Error(w, "Unknown fields are present in config-body. (Check for typos)", http.StatusBadRequest)
		log.Println("Unknown fields are present in attempted POST-body.", err)
		return nil, err
	}

	// Validated, return as []byte
	return bodyBytes, nil
}
