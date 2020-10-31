package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"montgomery.wg/blobreader/pkg/models"
	"montgomery.wg/blobreader/pkg/models/truservicelog"
)

func (app *truservicelog.Application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }


    reqs, err := app.httprequests.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    data := &templateData{TruServiceRequests: reqs}
    // Initialize a slice containing the paths to the show.page.tmpl file,
    // plus the base layout and footer partial that we made earlier.
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }


    // Parse the template files...
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // And then execute them. Notice how we are passing in the snippet
    // data (a models.Snippet struct) as the final parameter.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }

    // for _, httprequest := range reqs {
    //     fmt.Fprintf(w, "%v\n", httprequest)
    // }
}

func (app *truservicelog.Application) showTruserviceRequest(w http.ResponseWriter, r *http.Request) {
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

    data := &templateData{TruServiceRequest: req}
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

    // And then execute them. Notice how we are passing in the snippet
    // data (a models.Snippet struct) as the final parameter.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }
    // Write the snippet data as a plain-text HTTP response body.
    // fmt.Fprintf(w, "%v", req)
}



func (app *truservicelog.Application) healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello from Blobreader HTTP server!"))
}
