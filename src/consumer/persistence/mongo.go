package persistence

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

// MongoDB represents a wrapper around the mongodb connection
type MongoDB struct {
	db *mgo.Session
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
		db: db,
	}, nil
}

// CloseConnection closes current mongo db session
func (m *MongoDB) CloseConnection() {
	m.db.Close()
	log.Println("Closed database connection")
}

// LoadOptions read all the available options from the polls
func (m *MongoDB) LoadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}
