package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ladydascalie/earthquakes/bot"
	"github.com/ladydascalie/earthquakes/config"
	"github.com/ladydascalie/earthquakes/config/locales"
	"github.com/ladydascalie/earthquakes/db/mongo"
	"github.com/ladydascalie/earthquakes/templates"
	"github.com/ladydascalie/earthquakes/workers"
	"github.com/paulmach/go.geojson"
	"golang.org/x/crypto/acme/autocert"
)

var (
	plates string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	// setup the app config
	config.Setup()

	contents, err := ioutil.ReadFile("static/plates.json")
	if err != nil {
		log.Fatalln(err)
	}
	plates = string(contents)

	mongo.Begin()

	// Start the workers
	go workers.Run()
	go bot.Begin()
}

func main() {
	// todo: move endpoints out of main.go
	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.New()

	// setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	e.Static("static", "static")

	// setup api group
	api := e.Group("api/v1")
	api.GET("/feed", feedHandler)
	api.GET("/points", pointsHandler)

	// setup public group
	public := e.Group("")
	public.GET("/", listHandler)
	public.GET("/:locale/:id", singleHandler)
	public.GET("/display", displayHandler)

	if config.Env == "prod" {
		// setup pre router middleware
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("quakes.cable.fyi")
		e.AutoTLSManager.Cache = autocert.DirCache(".cache")
		e.StartAutoTLS(":443")
	} else {
		log.Fatalln(e.Start(config.AppPort))
	}
}

// this powers the feed api
func feedHandler(c echo.Context) error {
	alerts := mongo.GetLastHundred()
	return c.JSON(http.StatusOK, alerts)
}

// this retrieves the last 100 alerts and pushes them into a GeoJSON feature collection
// for consumption by the display template, generating the map.
func pointsHandler(c echo.Context) error {
	alerts := mongo.GetLastHundred()
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

// listHandler display
func listHandler(c echo.Context) error {
	alerts := mongo.GetLastHundred()
	return c.Render(200, "select", alerts)
}

// display handler serves the compiled display template
func displayHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "display", plates)
}

// singleHandler retrieves a single alert and returns
// the rendered `notify` template
func singleHandler(c echo.Context) error {
	if strings.HasPrefix(c.Request().RequestURI, "/api") {
		return c.NoContent(404)
	}
	id := c.Param("id")
	locale := c.Param("locale")
	locales.Load(locale)

	alert := mongo.GetById(id)

	return c.Render(200, "notify", alert)
}
