package main

import (
	"assigment2/handlers"
	"assigment2/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}
	http.HandleFunc(utils.ROOT_PATH, handlers.RootPath)

	log.Println("Starting server on port: " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
