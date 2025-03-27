package utils

import "time"

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
 *	Struct for dashboard configurations
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
	LastRetrieval time.Time `firestore:"lastChange" json:"lastChange"`
}

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
		Area             int                `json:"area,omitempty"`
		TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
	} `json:"features"`
	LastRetrieval time.Time `json:"lastRetrieval"`
}

/*
 *	Struct populated with the response from RestCountries
 */
type RestCountriesResponse struct {
	Capital     []string               `json:"capital"`
	Coordinates []int                  `json:"latlng"`
	Population  int                    `json:"population"`
	Area        int                    `json:"area"`
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
