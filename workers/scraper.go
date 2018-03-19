package workers

import (
	"fmt"
	"log"

	"github.com/paulmach/go.geojson"
	"github.com/ladydascalie/earthquakes/db/mongo"
	"github.com/ladydascalie/earthquakes/domain"
)

// Scrape is where the GeoJSON gets transformed into actual alerts models
func Scrape(_ int, jobs <-chan *geojson.Feature, results chan<- *domain.Alert) {
	for feature := range jobs {
		alert := new(domain.Alert)

		id := feature.ID.(string)
		if mongo.Exists(id) {
			return
		}

		detailsURL := feature.PropertyMustString("detail")
		details, err := getDetails(detailsURL)
		if err != nil {
			log.Println(err)
		}

		geoServeURL := getGeoServeURL(details)
		if geoServeURL == "" {
			return
		}

		geo, _ := getGeoServe(geoServeURL)
		if err != nil {
			log.Println(err)
		}
		for k, city := range geo.Cities {
			city.AppleMapsURL = fmt.Sprintf("http://maps.apple.com/?q=%s&ll=%f,%f&z=9", *city.Name, *city.Latitude, *city.Longitude)
			city.OpenStreetMapURL = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=8/%f/%f", *city.Latitude, *city.Longitude, *city.Latitude, *city.Longitude)
			geo.Cities[k] = city
		}

		// todo: url encode the apple maps url?
		// we use apple maps because they redirect to google
		// if you are not on an apple device.
		alert = &domain.Alert{
			ID:               id,
			Title:            feature.PropertyMustString("title"),
			Magnitude:        feature.PropertyMustFloat64("mag"),
			Place:            feature.PropertyMustString("place"),
			URL:              feature.PropertyMustString("url"),
			Created:          toTime(feature.PropertyMustFloat64("time")).UTC(),
			CreatedNano:      toTime(feature.PropertyMustFloat64("time")).UTC().UnixNano(),
			Lng:              feature.Geometry.Point[0],
			Lat:              feature.Geometry.Point[1],
			Depth:            feature.Geometry.Point[2],
			AffectedCities:   geo.Cities,
			AppleMapsURL:     fmt.Sprintf("http://maps.apple.com/?q=%s&ll=%f,%f&z=9", feature.PropertyMustString("title"), feature.Geometry.Point[1], feature.Geometry.Point[0]),
			OpenStreetMapURL: fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=8/%f/%f", feature.Geometry.Point[1], feature.Geometry.Point[0], feature.Geometry.Point[1], feature.Geometry.Point[0]),
		}

		// Sent the alert to the results channel
		results <- alert
		<-jobs // read from jobs
	}
}
