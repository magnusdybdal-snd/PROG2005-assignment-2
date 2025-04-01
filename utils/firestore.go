package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Client used across packages
var FirestoreClient *firestore.Client

func InitFirestore() error {

	// Firebase initialisation
	ctx := context.Background()

	opt := option.WithCredentialsFile("./api-keys/serviceAccountKey.json") // API KEY NEEDS TO BE LOCAL! - and added to .gitignore!
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
		return err
	}
	// Initiate client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Error initializing firestore: %v", err)
		return err
	}

	// Assign the initiated client to the global variable
	FirestoreClient = client
	return nil
}

// Close down client
func CloseFirestore() {
	if FirestoreClient != nil {
		errClose := FirestoreClient.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", errClose)
		}
	}
}