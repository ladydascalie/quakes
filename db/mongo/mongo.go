package mongo

import (
	"log"

	"github.com/ladydascalie/earthquakes/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Table in which the alert data will be stored
	Table = "alerts"
	// Collection under which the alerts will be stored
	Collection = "quakes"
)

// todo: move away from mgo and use community fork?

// Session is the global mongodb session object used for querying / inserting data
var Session *mgo.Session
// DB is the parent database object
var DB *mgo.Database

// Begin starts the connection to mongodb
func Begin() {
	var err error
	Session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	DB = Session.DB(Table)
}

// Exists check if any results at all exist for a given id
func Exists(id string) bool {
	if err := DB.C(Collection).Find(bson.M{"id": id}).Limit(1).One(&struct{}{}); err != nil {
		return false
	}
	return true
}

// GetById returns an alert by it's id
func GetById(id string) domain.Alert {
	var a domain.Alert
	if err := DB.C(Collection).Find(bson.M{"id": id}).One(&a); err != nil {
		log.Println(err)
		return a
	}
	return a
}

// GetLastHundred returns the last 100 alerts, sorted in descending order by their timestamp
func GetLastHundred() domain.Alerts {
	var alerts domain.Alerts
	DB.C(Collection).Find(nil).Sort("-creatednano").Limit(100).All(&alerts)
	return alerts
}
