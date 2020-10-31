package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"montgomery.wg/pedapilog/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }


    reqs, err := app.httprequests.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    data := &templateData{PEDAPIRequests: reqs}
    // Initialize a slice containing the paths to the show.page.tmpl file,
    // plus the base layout and footer partial that we made earlier.
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }

}

func (app *application) showPEDAPIRequest(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    req, err := app.httprequests.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    data := &templateData{PEDAPIRequest: req}
        // Initialize a slice containing the paths to the show.page.tmpl file,
    // plus the base layout and footer partial that we made earlier.
    files := []string{
        "./ui/html/show.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    // Parse the template files...
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
    if err := app.httprequests.DB.Ping(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }else{
        w.WriteHeader(http.StatusOK)
    }

}


func (app *application) readinessHandler(w http.ResponseWriter, r *http.Request) {
    if err := app.httprequests.DB.Ping(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }else{
        w.WriteHeader(http.StatusOK)
    }
}