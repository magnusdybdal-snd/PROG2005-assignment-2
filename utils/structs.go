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
