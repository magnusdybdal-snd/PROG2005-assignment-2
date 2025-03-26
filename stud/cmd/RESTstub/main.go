package main

import (
	"log"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8081")
		port = "8081"
	}

}