package config

import (
	"log"
	"os"
)

const (
	// also see:
	// https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson
	// https://earthquake.usgs.gov/earthquakes/eventpage/

	// USGSFeed is the default feed for significant seismic events
	USGSFeed = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_month.geojson"
)

var (
	// Env specifies the environment we are running in
	// this should be either dev or prod
	Env string

	// AppPort defines which port the application should run on
	AppPort string

	// BotToken defines the telegram bot token
	BotToken string

	// WithBot controls whether or not the bot will receive notifications
	WithBot bool
)

// Setup the application environment.
// this should be called shortly after loading the env file.
func Setup() {
	Env = MustGetEnv("ENV")
	AppPort = MustGetEnv("APP_PORT")
	BotToken = MustGetEnv("BOT_TOKEN")
	WithBot = false
}

// MustGetEnv will fatal if the key cannot be found
func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s key is missing from environment", key)
	}
	return val
}
