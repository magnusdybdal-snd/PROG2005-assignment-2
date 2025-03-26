package stubs

import (
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


/*
Responds with fixed JSON output sourced from provided file.
*/
func StubHandlerCountries(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		log.Println("Recieved " + r.Method + " request on Countries stub handler. Returning mocked information.")
		w.Header().Add("content-type", "application/json")
		output := ParseFile("./testdata/countries.json")
		fmt.Fprint(w, string(output))
		break
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}