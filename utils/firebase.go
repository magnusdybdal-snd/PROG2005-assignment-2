package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var ctx context.Context
var client *firestore.Client

func intializeFirebase() {
	ctx = context.Background()

	opt := option.WithCredentialsFile("../api-keys/serviceAccountKey.json") // API KEY NEEDS TO BE LOCAL! - and added to .gitignore!
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal("error initializing app: ", err)

	}

	client, err = app.Firestore(ctx)

	if err != nil {
		log.Fatal("Error initializing firestore: ", err)
	}

	defer client.Close()
}

func getFirebaseInfo() (context.Context, *firestore.Client) {
	if client != nil {
		intializeFirebase()
	}
	return ctx, client
}
