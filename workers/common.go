package workers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.com/ladydascalie/quakes/db/models"
	geojson "github.com/paulmach/go.geojson"
)

// Feed represents the data returned by querying the GeoJSON feed provided by the USGS
var Feed *geojson.FeatureCollection

// a LOT of type assertions need to be made here unfortunately
// however, safety checks are made so we can return an empty string rather than panic
func getGeoServeURL(details *geojson.Feature) (url string) {
	const jsonPath = "properties.products.geoserve.0.contents.geoserve\\.json.url"

	content, err := details.MarshalJSON()
	if err != nil {
		panic(err)
	}
	j := gjson.ParseBytes(content)
	if j.Get(jsonPath).Exists() {
		url = j.Get(jsonPath).String()
	}
	return
}

func getDetails(detailURL string) (*geojson.Feature, error) {
	res, err := http.Get(detailURL)
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
