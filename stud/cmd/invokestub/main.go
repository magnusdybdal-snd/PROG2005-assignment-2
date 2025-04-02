package main

import (
	"log"
	"net/http"
	"os"
	"stub/internal/stubs"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8082")
		port = "8082"
	}
	log.Println("Running on port: ", port)

	http.HandleFunc("/invoked", stubs.StubHandlerWebhook)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
