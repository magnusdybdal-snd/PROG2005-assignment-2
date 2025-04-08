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
		log.Println("$PORT has not been set. Default: 8081")
		port = "8081"
	}

	http.HandleFunc("/v3.1/alpha/", stubs.StubHandlerCountries)
	http.HandleFunc("/v1/forecast", stubs.StubHandlerWeather)
	http.HandleFunc("/currency/", stubs.StubHandlerCurrencies)
	http.HandleFunc("/invoked", stubs.StubHandlerWebhook)

	log.Println("Running on port: ", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
