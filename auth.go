package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var validKeys = []string{
	"open_sesame",
	"french^bread",
}

func authTest(fxn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		if strings.Contains(v["key"], "::") {
			log.Printf("INVALID REQUEST URL [%+v]\n", r.URL)
			http.Error(w, "Invalid Request", http.StatusTeapot)
			return
		}
		valid := false
		for _, key := range validKeys {
			if key == r.Header.Get("x-secret-squirrel") {
				valid = true
			}
			if key == r.URL.Query().Get("x-secret-squirrel") {
				valid = true
			}
		}
		if valid {
			fxn(w, r)
		} else {
			log.Printf("AUTHTEST FAIL [%+v]\n", r.URL)
			http.Error(w, "Not Authorized", http.StatusTeapot)
		}
	}
}
