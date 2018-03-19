package workers

import (
	"log"
	"os"
	"time"

	"github.com/paulmach/go.geojson"
	"github.com/ladydascalie/earthquakes/bot"
	"github.com/ladydascalie/earthquakes/db/mongo"
	"github.com/ladydascalie/earthquakes/domain"
)

var (
	jobs    = make(chan *geojson.Feature, 100)
	results = make(chan *domain.Alert, 100)
)

// Run kicks up the workers and sets them running
func Run() {
	for { // repeat the scheduled work at 15 second intervals
		feed := MonitorFeed()
		scheduledTasks(jobs, results, feed)
		time.Sleep(15 * time.Second)
	}
}

func scheduledTasks(jobs chan *geojson.Feature, results chan *domain.Alert, feed *geojson.FeatureCollection) {
	// Start by scraping
	go func() {
		for i := 0; i < 16; i++ { // spin 16 scrapers
			go Scrape(i, jobs, results)
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
			c := mongo.DB.C("quakes")
			if err := c.Insert(&alert); err != nil {
				log.Println(err)
				return
			}
			if os.Getenv("WITH_BOT") == "trues" {
				bot.NotifyTelegramChannel(alert)
			}

		}
	}()
}
