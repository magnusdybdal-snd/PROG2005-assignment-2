package main

import (
	"assignment2/handlers"
	"assignment2/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	// Initialize Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer utils.CloseFirestore()

	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)
	http.HandleFunc(utils.DASHBOARD_PATH, handlers.HandleGetDashboard)
	http.HandleFunc(utils.REGISTRATION_PATH+"{id}", handlers.HandleMessages)
	http.HandleFunc(utils.REGISTRATION_PATH, handlers.HandleMessages) // for get request without an {id}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
