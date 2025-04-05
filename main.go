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


	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8083")
		port = "8083"
	}
	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)
	http.HandleFunc(utils.DASHBOARD_PATH, handlers.HandleGetDashboard)
	http.HandleFunc(utils.NOTIFICATION_PATH+"{id}", handlers.HandleNotification)
	http.HandleFunc(utils.NOTIFICATION_PATH, handlers.HandleNotification)

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
