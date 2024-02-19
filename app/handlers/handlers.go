package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/ladydascalie/quakes/app"
	"github.com/ladydascalie/quakes/config/locales"
	"github.com/ladydascalie/quakes/db/mongo"
	geojson "github.com/paulmach/go.geojson"
)

var plates string

func init() {
	// read this file in memory and keep it there
	// todo: maybe use something like packr?
	contents, err := ioutil.ReadFile("static/plates.json")
	if err != nil {
		log.Fatalln(err)
	}
	plates = string(contents)
}

// FeedHandler this powers the feed api
func FeedHandler(server *app.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		alerts := mongo.Get(server.DB, 100)
		return c.JSON(http.StatusOK, alerts)
	}
}

// PointsHandler this retrieves the last 100 alerts and pushes them into a GeoJSON feature collection
// for consumption by the display template, generating the map.
func PointsHandler(server *app.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		alerts := mongo.Get(server.DB, 100)
		fc := geojson.NewFeatureCollection()

		for _, alert := range alerts {
			point := geojson.NewPointFeature([]float64{alert.Lng, alert.Lat})
			point.Properties = map[string]interface{}{
				"id":        alert.ID,
				"title":     alert.Title,
				"magnitude": alert.Magnitude,
				"depth":     alert.Depth,
				"date":      alert.Created.In(time.UTC).Format(time.RFC850),
			}
			fc.AddFeature(point)
		}
		return c.JSON(http.StatusOK, fc)
	}
}

// ListHandler display the list of quakes
func ListHandler(server *app.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var magnitude int
		var err error

		m := c.QueryParam("magnitude")
		magnitude, err = strconv.Atoi(m)
		if err != nil {
			magnitude = 0
		}

		alerts := mongo.GetMinimumMagnitude(server.DB, 100, magnitude)
		return c.Render(200, "select", alerts)
	}
}

// DisplayHandler serves mapbox display page
func DisplayHandler(_ *app.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "display", plates)
	}
}

// SingleHandler retrieves a single alert and returns
// the rendered `notify` template
func SingleHandler(server *app.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().RequestURI, "/api") {
			return c.NoContent(404)
		}

		id := c.Param("id")
		locale := c.Param("locale")
		locales.Load(locale)

		alert := mongo.GetById(server.DB, id)

		return c.Render(200, "notify", alert)
	}
}
