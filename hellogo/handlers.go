package main

import (
	"io/ioutil"
	"net/http"
)

func (app *application) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		w.Write([]byte("hello world"))
	} else {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	}

}
