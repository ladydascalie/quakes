package mongo

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ladydascalie/earthquakes/domain"
)

const (
	// Table in which the alert data will be stored
	Table = "alerts"
	// Collection under which the alerts will be stored
	Collection = "quakes"
)

// Begin starts the connection to mongodb
func Begin() *mgo.Database {
	var err error
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	return session.DB(Table)
}

// Exists check if any results at all exist for a given id
func Exists(db *mgo.Database, id string) bool {
	if err := db.C(Collection).Find(bson.M{"id": id}).Limit(1).One(&struct{}{}); err != nil {
		return false
	}
	return true
}

// GetById returns an alert by it's id
func GetById(db *mgo.Database, id string) domain.Alert {
	var a domain.Alert
	if err := db.C(Collection).Find(bson.M{"id": id}).One(&a); err != nil {
		log.Println(err)
		return a
	}
	return a
}

// Get returns the last 100 alerts, sorted in descending order by their timestamp
func Get(db *mgo.Database, limit int) domain.Alerts {
	var alerts domain.Alerts
	db.C(Collection).Find(nil).Sort("-creatednano").Limit(limit).All(&alerts)
	return alerts
}

// GetMinimumMagnitude the last 100 alerts, sorted in descending order by their timestamp
func GetMinimumMagnitude(db *mgo.Database, limit, magnitude int) domain.Alerts {
	var alerts domain.Alerts
	db.C(Collection).
		Find(bson.M{"magnitude": bson.M{"$gte": magnitude}}).
		Sort("-creatednano").
		Limit(limit).
		All(&alerts)
	return alerts
}
