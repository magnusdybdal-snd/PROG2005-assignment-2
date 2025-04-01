package utils

import (
	"context"
	"log"
)

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
func GetDashboardConfig[T any] (ctx context.Context, docId string, collection string) (T, error) {

	// Initiate zero value of type T
	var result T

	// Gets reference to firestore document
	docRef := FirestoreClient.Collection(collection).Doc(docId)

	// Fetches the data from the document
	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Println("Failed to get dashboard document %s: %v" + docId, err)
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