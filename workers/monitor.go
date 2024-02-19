package workers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ladydascalie/quakes/config"
	geojson "github.com/paulmach/go.geojson"
)

// MonitorFeed fetches the data from the USGS feed
func getFeed() *geojson.FeatureCollection {
	cli := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := cli.Get(config.USGSFeed)
	if err != nil {
		log.Println(err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	fc := new(geojson.FeatureCollection)
	if err := json.Unmarshal(b, fc); err != nil {
		log.Println(err)
		return nil
	}
	return fc
}
