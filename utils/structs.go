package utils

import "time"

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
	Id string
	DashboardConfig
}
