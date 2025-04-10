package utils

import "time"

// === Structs used in notification hanlder ===

type SendNotification struct {
	ID      string `json:"id"`
	Country string `firestore:"country" json:"country"`
	Event   string `firestore:"event" json:"event"`
	Time    string `json:"time"`
}
type RegisterWebhook struct {
	Url     string `firestore:"url" json:"url"`
	Country string `firestore:"country" json:"country"`
	Event   string `firestore:"event" json:"event"`
}

// When a get request is made, this is the struct to be returned
type ReturnWebhook struct {
	ID      string `json:"id"`
	Country string `firestore:"country" json:"country"`
	Event   string `firestore:"event" json:"event"`
	Url     string `firestore:"url" json:"url"`
}

// used to return the ID of the webhook
type WebhookId struct {
	ID string `json:"id"`
}


/*
 *	Struct for dashboard configurations, used by registratoin handler and dashboard handler
 */
type DashboardConfig struct {
	Country  string `firestore:"country" json:"country"`
	IsoCode  string `firestore:"isoCode" json:"isoCode"`
	Features struct {
		Temperature      bool     `firestore:"temperature" json:"temperature"`
		Precipitation    bool     `firestore:"precipitation" json:"precipitation"`
		Capital          bool     `firestore:"capital" json:"capital"`
		Coordinates      bool     `firestore:"coordinates" json:"coordinates"`
		Population       bool     `firestore:"population" json:"population"`
		Area             bool     `firestore:"area" json:"area"`
		TargetCurrencies []string `firestore:"targetCurrencies" json:"targetCurrencies"`
	} `firestore:"features" json:"features"`
	LastRetrieval string `firestore:"lastChange" json:"lastChange"`
}

// === Structs used in registrations handler ===

/*
 *	Struct for dashboard alterations (PUT) - without ID and TIME(!)
 */
type DashboardAlteration struct {
	Country  string `firestore:"country" json:"country"`
	IsoCode  string `firestore:"isoCode" json:"isoCode"`
	Features struct {
		Temperature      bool     `firestore:"temperature" json:"temperature"`
		Precipitation    bool     `firestore:"precipitation" json:"precipitation"`
		Capital          bool     `firestore:"capital" json:"capital"`
		Coordinates      bool     `firestore:"coordinates" json:"coordinates"`
		Population       bool     `firestore:"population" json:"population"`
		Area             bool     `firestore:"area" json:"area"`
		TargetCurrencies []string `firestore:"targetCurrencies" json:"targetCurrencies"`
	} `firestore:"features" json:"features"`
}

/*
* Helper struct to update timestamp for registration PUT function.
 */
type DashboardAlterationTime struct {
	DashboardAlteration
	LastRetrieval string `firestore:"lastChange" json:"lastChange"`
}


// === Structs used in dahsboard handler ===

/*
 *	Response struct from dashboards handler
 */
type DashboardResponse struct {
	Country  string `json:"country"`
	IsoCode  string `json:"isoCode"`
	Features struct {
		Temperature      float64            `json:"temperature,omitempty"`
		Precipitation    float64            `json:"precipitation,omitempty"`
		Capital          string             `json:"capital,omitempty"`
		Coordinates      map[string]float64 `json:"coordinates,omitempty"`
		Population       int                `json:"population,omitempty"`
		Area             float64            `json:"area,omitempty"`
		TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
	} `json:"features"`
	LastRetrieval string `json:"lastRetrieval"`
}

/*
 *	Struct populated with the response from RestCountries
 */
type RestCountriesResponse struct {
	Capital     []string               `json:"capital"`
	Coordinates []float64              `json:"latlng"`
	Population  int                    `json:"population"`
	Area        float64                `json:"area"`
	Currencies  map[string]interface{} `json:"currencies"`
}

/*
 *	Struct populated with the response from MetroAPI
 */
type MetroResponse struct {
	Hourly struct {
		Precipitation []float64 `json:"precipitation"`
		Temperature2M []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

/*
*	Struct returned from getMetroData()
 */
type MetroMeanValues struct {
	MeanPrecipitation float64
	MeanTemperature   float64
}

/*
 *	Struct populated with the response from CurrencyAPI
 */
type CurrencyResponse struct {
	Rates map[string]float64 `json:"rates"`
}

/*
* Struct for registration GET responses
* ..using DashboardConfig
 */
type RegistrationGetResponse struct {
	Id string `json:"id"`
	DashboardConfig
}


/*
*	Struct used to hold cached information from RestCountries
*/
type CachedRestCountries struct {
	Data      RestCountriesResponse `firestore:"data"`
	Timestamp time.Time             `firestore:"timestamp"`
}

/*
*	Struct used to hold cached information from Metro API
*/
type CachedMetro struct {
	Data      MetroMeanValues `firestore:"data"`
	Timestamp time.Time       `firestore:"timestamp"`
}

/*
*	Struct used to hold cached information from Currency API
*/
type CachedCurrency struct {
	Data      map[string]float64 `firestore:"data"`
	Timestamp time.Time          `firestore:"timestamp"`
}

/*
*	Struct used to hold status information in status endpoint
*/
type Status struct {
	Countries_api   int    `json:"countries_api"`
	Metro_api       int    `json:"metro_api"`
	Currency_api    int    `json:"currency_api"`
	Notification_db int    `json:"notification_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}