package utils

import (
	"bytes"
	"context"
	"encoding/json"
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
