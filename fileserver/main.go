// package main

// import (
// 	"net/http"
// 	"os"

// 	"github.com/russross/blackfriday"
// )

// func main() {

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

//     http.HandleFunc("/markdown", GenerateMarkdown)
//     http.Handle("/", http.FileServer(http.Dir("public")))
//     http.ListenAndServe(":"+port, nil)
// }

// func GenerateMarkdown(rw http.ResponseWriter, r *http.Request) {
//     markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
//     rw.Write(markdown)
// }

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

    d := (24 * time.Hour) 

    // Calling Round() method 
    then := time.Now().UTC().Add(-24 * time.Hour).Round(d).Format("2006-01-02")
    now := time.Now().UTC().Round(d).Format("2006-01-02")

    fmt.Printf("Between %s and %s\n", then, now)

    r := httprouter.New()
    r.GET("/", HomeHandler)

    // Posts collection
    r.GET("/posts", PostsIndexHandler)
    r.POST("/posts", PostsCreateHandler)

    // Posts singular
    r.GET("/posts/:id", PostShowHandler)
    r.PUT("/posts/:id", PostUpdateHandler)
    r.GET("/posts/:id/edit", PostEditHandler)

    

    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintln(rw, "Home")
}

func PostsIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    t := time.Now().UTC() //.Add(-5 * time.Hour)
    // Defining duration 
    d := (60 * time.Minute) 

    // Calling Round() method 
    res := t.Round(d).Format("2006-01-02 15:04")
    fmt.Fprintln(rw, res)
}

func PostsCreateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintln(rw, "posts create")
}

func PostShowHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    id := p.ByName("id")
    fmt.Fprintln(rw, "showing post", id)
}

func PostUpdateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintln(rw, "post update")
}

func PostDeleteHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintln(rw, "post delete")
}

func PostEditHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintln(rw, "post edit")
}