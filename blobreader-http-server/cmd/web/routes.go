package main

import (
	"net/http"
	"path/filepath"

	"montgomery.wg/blobreader/pkg/models/truservicelog"
)


func (app *truservicelog.Application) routes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/request", app.showTruserviceRequest)
    mux.HandleFunc("/healthz", app.healthHandler)
	mux.HandleFunc("/ready", app.healthHandler)

    fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
    mux.Handle("/static", http.NotFoundHandler())
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
} 