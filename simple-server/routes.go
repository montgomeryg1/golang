package main

import (
	"github.com/gorilla/mux"
)


func (app *application) routes() *mux.Router {
    //mux := http.NewServeMux()
    
    r := mux.NewRouter()
	// Route handles & endpoints
    r.HandleFunc("/check", app.check).Methods("GET")

    r.HandleFunc("/addBook", app.addBook).Methods("GET")
    
    // mux.HandleFunc("/", app.home)
    // mux.HandleFunc("/snippet", app.showSnippet)
    // mux.HandleFunc("/snippet/create", app.createSnippet)

    // fileServer := http.FileServer(http.Dir("./ui/static/"))
    // mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    return r
}