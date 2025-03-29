package utils

import "time"

/*
*	Struct for dashboard configuration.
 */
type DashboardConfig struct {
	Country string					`firestore:"country" json:"country"`
	IsoCode string					`firestore:"isoCode" json:"isoCode"`
	Features struct {
		Temperature bool			`firestore:"temperature" json:"temperature"`
		Precipiation bool			`firestore:"percipiation" json:"precipiation"`
		Capital bool				`firestore:"capital" json:"capital"`
		Coordinates bool			`firestore:"coordinates" json:"coordinates"`
		Population bool				`firestore:"population" json:"population"`
		Area bool					`firestore:"area" json:"area"`
		TargetCurrencies []string 	`firestore:"targetCurrencies" json:"targetCurrencies"`
	}								`firestore:"features" json:"features"`
	LastRetrieval time.Time			`firestore:"lastChange" json:"lastChange"`
}


/*
*	Struct populated with the response from RestCountries
*/
type RestCountriesResponse struct {
	Capital []string 			   	`json:"capital"`
	Coordinates []int		       	`json:"latlng"`
	Population int                 	`json:"population"`
	Area int						`json:"area"`
}

/*
*	Struct populated with the response from MetroAPI
*/
type MetroResponse struct {
	Hourly struct {
		Precipitation []float64     `json:"precipitation"`
		Temperature2M []float64    	`json:"temperature_2m"`
	} 								`json:"hourly"`
}

/*
*	Struct populated with the response from CurrencyAPI
*/
type CurrencyResponse struct {
	Rates map[string]float64		`json:"rates"`
}