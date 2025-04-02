package utils

import (
	"assignment2/utils"
	"context"
	"log"
)

func CallUrl(url string, event string, content utils.SendNotification)

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
