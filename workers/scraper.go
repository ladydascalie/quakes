package workers

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/ladydascalie/quakes/db/mongo"
	"github.com/ladydascalie/quakes/domain"
	geojson "github.com/paulmach/go.geojson"
)

// Scrape is where the GeoJSON gets transformed into actual alerts models
func processAlerts(db *mgo.Database, jobs <-chan *geojson.Feature, results chan<- *domain.Alert) {
	for feature := range jobs {
		// ensure entry is cleaned
		defer func() { <-jobs }()

		id, ok := feature.ID.(string)
		if !ok {
			log.Println("unexpected: feature id is not a string")
			return
		}

		if mongo.Exists(db, id) {
			return
		}

		// detailsURL := feature.PropertyMustString("detail")
		// details, err := getDetails(detailsURL)
		// if err != nil {
		// 	log.Println(err)
		// }

		// geoServeURL := getGeoServeURL(details)
		// if geoServeURL == "" {
		// 	log.Println("no geoserve url")
		// 	return
		// }

		// geoServe, err := getGeoServe(geoServeURL)
		// if err != nil {
		// 	log.Println(err)
		// }
		// for k, city := range geoServe.Cities {
		// 	city.AppleMapsURL = fmt.Sprintf("http://maps.apple.com/?q=%s&ll=%f,%f&z=9", *city.Name, *city.Latitude, *city.Longitude)
		// 	city.OpenStreetMapURL = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=8/%f/%f", *city.Latitude, *city.Longitude, *city.Latitude, *city.Longitude)
		// 	geoServe.Cities[k] = city
		// }

		alert := alertFromFeature(feature)
		if alert == nil {
			log.Println("alertFromFeature: returned nil")
			return
		}

		// Sent the alert to the results channel
		results <- alert
	}
}

func alertFromFeature(feature *geojson.Feature) *domain.Alert {
	id, ok := feature.ID.(string)
	if !ok {
		return nil
	}
	// todo: url encode the apple maps url?
	// we use apple maps because they redirect to google
	// if you are not on an apple device.
	return &domain.Alert{
		ID:          id,
		Title:       feature.PropertyMustString("title"),
		Magnitude:   feature.PropertyMustFloat64("mag"),
		Place:       feature.PropertyMustString("place"),
		URL:         feature.PropertyMustString("url"),
		Created:     toTime(feature.PropertyMustFloat64("time")).UTC(),
		CreatedNano: toTime(feature.PropertyMustFloat64("time")).UTC().UnixNano(),
		Lng:         feature.Geometry.Point[0],
		Lat:         feature.Geometry.Point[1],
		Depth:       feature.Geometry.Point[2],
		// AffectedCities:   geoServe.Cities,
		AppleMapsURL:     fmt.Sprintf("http://maps.apple.com/?q=%s&ll=%f,%f&z=9", feature.PropertyMustString("title"), feature.Geometry.Point[1], feature.Geometry.Point[0]),
		OpenStreetMapURL: fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=8/%f/%f", feature.Geometry.Point[1], feature.Geometry.Point[0], feature.Geometry.Point[1], feature.Geometry.Point[0]),
	}
}
