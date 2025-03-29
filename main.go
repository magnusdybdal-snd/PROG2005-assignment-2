package main

import (
	"assignment2/handlers"
	"assignment2/utils"
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var ctx context.Context
var client *firestore.Client

func main() {

	ctx = context.Background()

	opt := option.WithCredentialsFile("/mnt/c/Users/magnu/OneDrive - NTNU/Skrivebord/NTNU/4SEM/PROG2005/assignment-2/serviceAccountKey.json") // API KEY NEEDS TO BE LOCAL! - and added to .gitignore!
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
		return
	}

	client, err = app.Firestore(ctx)

	if err != nil {
		log.Printf("Error initializing firestore: %v", err)
		return
	}

	defer client.Close()

	dashboardHandler := handlers.NewDashboardHandler(client)

	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)
	http.HandleFunc(utils.DASHBOARD_PATH, dashboardHandler.ServeHTTP)

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
