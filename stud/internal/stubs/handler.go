package stubs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

/*
Reads a given file and returns content as byte array.
*/
func ParseFile(filename string) []byte {
	file, e := os.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file
}

func StubHandlerWebhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var webhook Webhook
		log.Println("Received " + r.Method + " request on invoke stub handler")

		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println("Bad request, wasn't able to decode JSON: ", err)
			return
		}

		log.Println("Received payload, sending result")
		log.Println(webhook)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(webhook); err != nil {
			log.Println("Unable to send payload: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
	}
}

/*
Responds with fixed JSON output sourced from provided file.
*/
func StubHandlerCountries(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		log.Println("Recieved " + r.Method + " request on Countries stub handler. Returning mocked information.")
		w.Header().Add("content-type", "application/json")
		output := ParseFile("../../testdata/countries.json")
		fmt.Fprint(w, string(output))
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

/*
Responds with fixed JSON output sourced from provided file.
*/
func StubHandlerCurrencies(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		log.Println("Recieved " + r.Method + " request on Currencies stub handler. Returning mocked information.")
		w.Header().Add("content-type", "application/json")
		output := ParseFile("../../testdata/currencies.json")
		fmt.Fprint(w, string(output))
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

/*
Responds with fixed JSON output sourced from provided file.
*/
func StubHandlerWeather(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		log.Println("Recieved " + r.Method + " request on Weather stub handler. Returning mocked information.")
		w.Header().Add("content-type", "application/json")
		output := ParseFile("../../testdata/weather.json")
		fmt.Fprint(w, string(output))
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}
