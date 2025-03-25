package main

import (
	"log"
	"net/http"
	"os"
)

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port has not been set, using default 8080")
		port = "8080"
	}

	log.Println("Starting server on port: "+ port +"..." )
	log.Fatal(http.ListenAndServe(":" + port, nil))
}