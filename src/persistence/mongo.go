package persistence

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

type poll struct {
	Options []string
}

// DialDB opens a connection to the database
func DialDB() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("twitter_votes_database")
	return err
}

// CloseDB closes the current connection to the database
func CloseDB() {
	db.Close()
	log.Println("Closed database connection")
}

// LoadOptions read all the available options from the polls
func LoadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}
