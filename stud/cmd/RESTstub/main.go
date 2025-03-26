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

	http.HandleFunc("/countries/no", stubs.StubHandlerCountries)
	http.HandleFunc("/weather/no", stubs.StubHandlerWeather)
	http.HandleFunc("/currency/no", stubs.StubHandlerCurrencies)

	log.Println("Running on port: ", port)

	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}