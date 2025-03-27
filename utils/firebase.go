package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Ctx context.Context
var Client *firestore.Client

func IntializeFirebase() {
	Ctx = context.Background()
	//todo: fix the file location
	opt := option.WithCredentialsFile("/home/olemgl/Documents/Skole/sem4/assignment2/api-keys/serviceAccountKey.json") // API KEY NEEDS TO BE LOCAL! - and added to .gitignore!

	app, err := firebase.NewApp(Ctx, nil, opt)
	if err != nil {
		log.Fatal("error initializing app: ", err)
	}

	Client, err = app.Firestore(Ctx)

	if err != nil {
		log.Fatal("Error initializing firestore: ", err)
	}
}

func CloseFirebase() {
	if Client != nil {
		Client.Close()
	}
}
