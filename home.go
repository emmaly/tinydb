package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupHomeMethods(r *mux.Router) {
	r.HandleFunc("/", authTest(homeHandler))
	r.HandleFunc("/favicon.ico", notFound)
	r.HandleFunc("/robots.txt", notFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := "homeHandler"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}
