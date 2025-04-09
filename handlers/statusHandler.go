package handlers

import (
	"assignment2/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
)

var startTime time.Time

func InitStart(t time.Time) {
	startTime = t
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported, please use "+http.MethodGet, http.StatusMethodNotAllowed)
		return
	}
	MetroUrl := fmt.Sprintf(utils.MetroAPI, 52.52, 13.41)
	status := utils.Status{
		Countries_api:   getHttpCode(utils.RESTCountriesAPI+"no", r),
		Metro_api:       getHttpCode(MetroUrl, r),
		Currency_api:    getHttpCode(utils.CurrencyAPI+"NOK", r),
		Notification_db: getFireStoreCode(r),
		Webhooks:        getAmmountWebhooks(r),
		Version:         utils.VERSION,
		Uptime:          int(time.Since(startTime).Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding status page: ", err)
		return
	}
}

func getAmmountWebhooks(r *http.Request) int {
	webhooks, err := retriveWebhooks(r)
	if err != nil {
		log.Println("Error fetching webhooks: ", err)
		return 0
	}
	return len(webhooks)
}

func getFireStoreCode(r *http.Request) int {
	client := utils.FirestoreClient
	iter := client.Collection(utils.WebhooksCollection).Limit(1).Documents(context.Background())
	defer iter.Stop()
	_, err := iter.Next()
	if err == iterator.Done {
		log.Println("No document in Notifications_DB: ", err)
		return 404
	} else if err != nil {
		log.Println("Error fetching document: ", err)
		invokeWebhook(utils.ACCESS_FAILURE, "", r)
		return 500

	}
	return 200
}

func getHttpCode(url string, r *http.Request) int {
	client := http.Client{}
	defer client.CloseIdleConnections()

	resp, err := client.Get(url)
	if err != nil {
		invokeWebhook(utils.NOTREACHABLE, "", r)
		log.Println("Unable to get API: "+url+" err: ", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
