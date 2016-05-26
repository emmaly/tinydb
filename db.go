package main

import (
	"log"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func setupDB() {
	var err error
	db, err = bolt.Open("tiny.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}
