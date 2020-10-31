package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"montgomery.wg/blobreader/pkg/models/truservicelog"
)


func (app *truservicelog.Application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Output(2, trace)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *truservicelog.Application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}


func (app *truservicelog.Application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}