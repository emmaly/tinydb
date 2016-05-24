package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var validKeys = []string{
	"open_sesame",
	"french^bread",
}

var db *bolt.DB

func main() {
	var err error
	db, err = bolt.Open("tiny.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", authTest(homeHandler))
	s := r.Path("/data/{bucket}/{key}").Subrouter()
	s.Methods("GET").HandlerFunc(authTest(keyGet))
	s.Methods("POST", "PUT").HandlerFunc(authTest(keyPost))
	s.Methods("DELETE").HandlerFunc(authTest(keyDelete))
	s.Methods("HEAD").HandlerFunc(authTest(keyHead))
	s.NewRoute().HandlerFunc(authTest(keyDefault))
	r.PathPrefix("/").HandlerFunc(authTest(notFound))
	http.Handle("/", r)
	http.ListenAndServe("0.0.0.0:8123", nil)
}

func authTest(fxn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.Println("AUTHTEST FAIL")
			http.Error(w, "Not Authorized", http.StatusTeapot)
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := "homeHandler"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func keyGet(w http.ResponseWriter, r *http.Request) {
	t := "keyGet"
	v := mux.Vars(r)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v["bucket"]))
		mt := new(time.Time)
		err := mt.UnmarshalBinary(b.Get([]byte(v["key"] + "/modified")))
		if err != nil {
			return err
		}
		w.Header().Add("date", mt.Format(time.RFC1123))
		d := b.Get([]byte(v["key"]))
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else if ct := b.Get([]byte(v["key"] + "/content-type")); len(ct) != 0 {
			contentType = string(ct)
		} else {
			contentType = http.DetectContentType(d)
		}
		w.Header().Add("content-type", contentType)
		w.Write(d)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
}

func keyPost(w http.ResponseWriter, r *http.Request) {
	t := "keyPost"
	v := mux.Vars(r)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(v["bucket"]))
		if err != nil {
			return err
		}
		defer r.Body.Close()
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		t, _ := time.Now().MarshalBinary()
		b.Put([]byte(v["key"]+"/modified"), t)
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else {
			contentType = http.DetectContentType(d)
		}
		b.Put([]byte(v["key"]+"/length"), []byte(fmt.Sprintf("%d", len(d))))
		b.Put([]byte(v["key"]+"/content-type"), []byte(contentType))
		b.Put([]byte(v["key"]), d)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	http.Error(w, "", http.StatusOK)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func keyDelete(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(v["bucket"]))
		if err != nil {
			return err
		}
		b.Delete([]byte(v["key"]))
		b.Delete([]byte(v["key"] + "/length"))
		b.Delete([]byte(v["key"] + "/content-type"))
		b.Delete([]byte(v["key"] + "/modified"))
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	http.Error(w, "", http.StatusOK)
}

func keyHead(w http.ResponseWriter, r *http.Request) {
	t := "keyHead"
	v := mux.Vars(r)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v["bucket"]))
		mt := new(time.Time)
		err := mt.UnmarshalBinary(b.Get([]byte(v["key"] + "/modified")))
		if err != nil {
			return err
		}
		w.Header().Add("date", mt.Format(time.RFC1123))
		d := b.Get([]byte(v["key"]))
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else if ct := b.Get([]byte(v["key"] + "/content-type")); len(ct) != 0 {
			contentType = string(ct)
		} else {
			contentType = http.DetectContentType(d)
		}
		w.Header().Add("content-type", contentType)
		if l := b.Get([]byte(v["key"] + "/length")); len(l) != 0 {
			w.Header().Add("content-length", string(l))
		} else {
			w.Header().Add("content-length", string(len(b.Get([]byte(v["key"])))))
		}
		return nil
	})
}

func keyDefault(w http.ResponseWriter, r *http.Request) {
	t := "keyDefault"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusMethodNotAllowed)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}
