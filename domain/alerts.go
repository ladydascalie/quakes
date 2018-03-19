package domain

import (
	"time"

	"github.com/ladydascalie/earthquakes/db/models"
)

// Alerts is a slice of Alert types, mostly used to contain a result set.
type Alerts []Alert

// Alert defines all the elements necessary to produce and store an alert.
type Alert struct {
	ID               string        `json:"id"`
	Title            string        `json:"title"`
	Magnitude        float64       `json:"magnitude"`
	Place            string        `json:"place"`
	URL              string        `json:"url"`
	Created          time.Time     `json:"created"`
	CreatedNano      int64         `json:"created_nano"`
	Lng              float64       `json:"lng"`
	Lat              float64       `json:"lat"`
	Depth            float64       `json:"depth"`
	AffectedCities   []models.City `json:"affectedCities"`
	OpenStreetMapURL string        `json:"open_street_map"`
	AppleMapsURL     string        `json:"apple_maps"`
}
