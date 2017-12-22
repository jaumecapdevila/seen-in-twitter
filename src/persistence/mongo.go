package persistence

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

// MongoDB represents a wrapper around the mongodb connection
type MongoDB struct {
	DB *mgo.Session
}
type poll struct {
	Options []string
}

// New returns a new MongoDB connection or an error otherwise
func New(source string) (*MongoDB, error) {
	var err error
	db, err = mgo.Dial(source)
	if err != nil {
		return &MongoDB{}, err
	}
	return &MongoDB{
		DB: db,
	}, nil
}

func (m *MongoDB) closeConnection() {
	m.DB.Close()
	log.Println("Closed database connection")
}

func (m *MongoDB) loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}
