package utils

import (
	"context"
	"log"
)

func GetDashboardConfig (ctx context.Context, dashboardId string, collection string) (DashboardConfig, error) {

	// Gets reference to firestore document
	docRef := FirestoreClient.Collection(collection).Doc(dashboardId)

	// Fetches the data from the document
	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Println("Failed to get dashboard document %s: %v" + dashboardId, err)
		return DashboardConfig{}, err
	}

	// Unmarshals the data in the document to struct
	var config DashboardConfig
	err2 := doc.DataTo(&config)
	if err2 != nil {
		log.Printf("Failed to unmarshal dashboard config: %s: %v", dashboardId, err)
		return DashboardConfig{}, err
	}

	return config, nil
}