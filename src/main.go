package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

func main() {
	fmt.Println("Hello world")
}

func dialDB() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("twitter_votes_database")
	return err
}

func closeDB() {
	db.Close()
	log.Println("Closed database connection")
}
