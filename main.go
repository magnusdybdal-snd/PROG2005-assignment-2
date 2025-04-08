package main

import (
	"assignment2/handlers"
	"assignment2/utils"
	"log"
	"net/http"
	"os"
	"time"
)

var startTime time.Time

// initate start time
func init() {
	startTime = time.Now()
}

func main() {
	//set time
	handlers.InitStart(startTime)
	// Initialize Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer utils.CloseFirestore()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}
	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)
	http.HandleFunc(utils.STATUS_PATH, handlers.StatusHandler)
	http.HandleFunc(utils.DASHBOARD_PATH, handlers.HandleGetDashboard)
	http.HandleFunc(utils.NOTIFICATION_PATH+"{id}", handlers.HandleNotification)
	http.HandleFunc(utils.NOTIFICATION_PATH, handlers.HandleNotification)
	http.HandleFunc(utils.REGISTRATION_PATH+"{id}", handlers.HandleMessages)
	http.HandleFunc(utils.REGISTRATION_PATH, handlers.HandleMessages) // for get request without an {id}

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
