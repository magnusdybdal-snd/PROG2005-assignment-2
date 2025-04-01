package main

import (
	"assignment2/handlers"
	"assignment2/utils"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {

	// Initialize Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer utils.CloseFirestore()

	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)
	http.HandleFunc(utils.DASHBOARD_PATH, handlers.HandleGetDashboard)

	// Pass the client to the handlers
	handlers.InitFirestore(client, ctx)

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}