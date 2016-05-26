package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupSubMethods(r *mux.Router) {
	s := r.PathPrefix("/sub/{bucket}/{key:.+}").Subrouter()
	s.Methods("GET").HandlerFunc(authTest(subGet))
	s.Methods("POST", "PUT").HandlerFunc(authTest(subPost))
	s.Methods("DELETE").HandlerFunc(authTest(subDelete))
	s.Methods("HEAD").HandlerFunc(authTest(subHead))
	s.NewRoute().HandlerFunc(authTest(subDefault))
}

func subGet(w http.ResponseWriter, r *http.Request) {
	t := "subGet"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func subPost(w http.ResponseWriter, r *http.Request) {
	t := "subPost"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func subDelete(w http.ResponseWriter, r *http.Request) {
	t := "subDelete"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func subHead(w http.ResponseWriter, r *http.Request) {
	t := "subHead"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func subDefault(w http.ResponseWriter, r *http.Request) {
	t := "subDefault"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}
