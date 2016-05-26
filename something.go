package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// ErrNoBucket means the requested bucket doesn't exist
var ErrNoBucket = errors.New("bucket doesn't exist")

func main() {
	setupDB()
	defer db.Close()

	r := mux.NewRouter()
	http.Handle("/", r)

	setupHomeMethods(r)
	setupDataMethods(r)
	setupSubMethods(r)

	r.PathPrefix("/").HandlerFunc(authTest(notFound))
	http.ListenAndServe("0.0.0.0:8123", nil)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
