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

func setupDataMethods(r *mux.Router) {
	dh := r.PathPrefix("/data/{bucket}/{key:.+}").Subrouter()
	dh.Methods("GET").HandlerFunc(authTest(dataGet))
	dh.Methods("POST", "PUT").HandlerFunc(authTest(dataPost))
	dh.Methods("DELETE").HandlerFunc(authTest(dataDelete))
	dh.Methods("HEAD").HandlerFunc(authTest(dataHead))
	dh.NewRoute().HandlerFunc(authTest(dataDefault))
}

func dataGet(w http.ResponseWriter, r *http.Request) {
	t := "keyGet"
	v := mux.Vars(r)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v["bucket"]))
		if b == nil {
			return ErrNoBucket
		}
		mt := time.Time{}
		mod := b.Get([]byte(v["key"] + "::modified"))
		if mod != nil {
			err := mt.UnmarshalBinary(mod)
			if err != nil {
				return err
			}
			w.Header().Add("date", mt.Format(time.RFC1123))
		}
		d := b.Get([]byte(v["key"]))
		if d == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else if ct := b.Get([]byte(v["key"] + "::content-type")); len(ct) != 0 {
			contentType = string(ct)
		} else {
			contentType = http.DetectContentType(d)
		}
		w.Header().Add("content-type", contentType)
		w.Write(d)
		return nil
	})
	if err != nil {
		if err == ErrNoBucket {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
}

func dataPost(w http.ResponseWriter, r *http.Request) {
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
		b.Put([]byte(v["key"]+"::modified"), t)
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else {
			contentType = http.DetectContentType(d)
		}
		b.Put([]byte(v["key"]+"::length"), []byte(fmt.Sprintf("%d", len(d))))
		b.Put([]byte(v["key"]+"::content-type"), []byte(contentType))
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

func dataDelete(w http.ResponseWriter, r *http.Request) {
	t := "keyDelete"
	v := mux.Vars(r)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(v["bucket"]))
		if err != nil {
			return err
		}
		b.Delete([]byte(v["key"]))
		b.Delete([]byte(v["key"] + "::length"))
		b.Delete([]byte(v["key"] + "::content-type"))
		b.Delete([]byte(v["key"] + "::modified"))
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}

func dataHead(w http.ResponseWriter, r *http.Request) {
	t := "keyHead"
	v := mux.Vars(r)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(v["bucket"]))
		if b == nil {
			return ErrNoBucket
		}
		mt := time.Time{}
		mod := b.Get([]byte(v["key"] + "::modified"))
		if mod != nil {
			err := mt.UnmarshalBinary(mod)
			if err != nil {
				return err
			}
			w.Header().Add("date", mt.Format(time.RFC1123))
		}
		d := b.Get([]byte(v["key"]))
		if d == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return nil
		}
		var contentType string
		if ct := r.Header.Get("x-content-type"); ct != "" {
			contentType = ct
		} else if ct := b.Get([]byte(v["key"] + "::content-type")); len(ct) != 0 {
			contentType = string(ct)
		} else {
			contentType = http.DetectContentType(d)
		}
		w.Header().Add("content-type", contentType)
		if l := b.Get([]byte(v["key"] + "::length")); len(l) != 0 {
			w.Header().Add("content-length", string(l))
		} else {
			w.Header().Add("content-length", string(len(b.Get([]byte(v["key"])))))
		}
		//w.Write(d)
		return nil
	})
	if err != nil {
		if err == ErrNoBucket {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
}

func dataDefault(w http.ResponseWriter, r *http.Request) {
	t := "keyDefault"
	v := mux.Vars(r)
	http.Error(w, "", http.StatusMethodNotAllowed)
	log.Printf("[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
	fmt.Fprintf(w, "[%s:%s] [%s/%s]\n", t, r.Method, v["bucket"], v["key"])
}
