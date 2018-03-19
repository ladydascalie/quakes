package workers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/paulmach/go.geojson"
	"github.com/ladydascalie/earthquakes/db/models"
)

// Feed represents the data returned by querying the GeoJSON feed provided by the USGS
var Feed *geojson.FeatureCollection

// a LOT of type assertions need to be made here unfortunately
// however, safety checks are made so we can return an empty string rather than panic
func getGeoServeURL(details *geojson.Feature) string {
	products, ok := details.Properties["products"].(map[string]interface{})
	if !ok {
		return ""
	}
	geoserve, ok := products["geoserve"].([]interface{})
	if !ok {
		return ""
	}
	zero, ok := geoserve[0].(map[string]interface{})
	if !ok {
		return ""
	}

	contents, ok := zero["contents"].(map[string]interface{})
	if !ok {
		return ""
	}
	geoserveJson, ok := contents["geoserve.json"].(map[string]interface{})
	if !ok {
		return ""
	}
	url, ok := geoserveJson["url"].(string)
	if !ok {
		return ""
	}

	return url
}

func getDetails(detailUrl string) (*geojson.Feature, error) {
	res, err := http.Get(detailUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	feature, err := geojson.UnmarshalFeature(body)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func getGeoServe(geoServeURL string) (*models.GeoServe, error) {
	res, err := http.Get(geoServeURL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	geoServe := new(models.GeoServe)
	json.Unmarshal(body, geoServe)
	if err != nil {
		return nil, err
	}
	return geoServe, nil
}

func toTime(timestamp float64) time.Time {
	// the timestamps are given in milliseconds
	i := int64(timestamp) / 1000
	tm := time.Unix(i, 0).UTC()

	return tm
}
