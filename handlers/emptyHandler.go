package handlers

import (
	"assignment2/utils"
	"fmt"
	"log"
	"net/http"
)

func RootPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")

	output := "This service provides functionality on these paths: <br>" +
		"<a href=\"" + utils.REGISTRATION_PATH + "\">" + utils.REGISTRATION_PATH + "<br>" +
		"<a href=\"" + utils.DASHBOARD_PATH + "\">" + utils.DASHBOARD_PATH + "<br>" +
		"<a href=\"" + utils.NOTIFICATION_PATH + "\">" + utils.NOTIFICATION_PATH + "<br>" +
		"<a href=\"" + utils.STATUS_PATH + "\">" + utils.STATUS_PATH

	_, err := fmt.Fprint(w, output)
	if err != nil {
		log.Println("Error when returning root output: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
