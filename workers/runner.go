package workers

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/ladydascalie/quakes/bot"
	"github.com/ladydascalie/quakes/config"
	"github.com/ladydascalie/quakes/domain"
	geojson "github.com/paulmach/go.geojson"
)

var (
	jobs    = make(chan *geojson.Feature, 100)
	results = make(chan *domain.Alert, 100)
)

// Run kicks up the workers and sets them running
func Run(db *mgo.Database) {
	// Do it once on startup
	process(db, jobs, results, getFeed())

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		process(db, jobs, results, getFeed())
	}
}

func process(db *mgo.Database, jobs chan *geojson.Feature, results chan *domain.Alert, feed *geojson.FeatureCollection) {
	// Start by scraping
	go func() {
		for i := 0; i < 16; i++ { // spin 16 scrapers
			go processAlerts(db, jobs, results)
		}
	}()

	go func() {
		for _, j := range feed.Features {
			jobs <- j // push onto the job channel
		}
	}()

	go func() {
		for range results {
			alert := <-results
			c := db.C("quakes")
			if err := c.Insert(&alert); err != nil {
				log.Println(err)
				return
			}
			if config.WithBot {
				bot.NotifyTelegramChannel(alert)
			}
		}
	}()
}
