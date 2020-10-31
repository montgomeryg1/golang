package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"montgomery.wg/http-server/pkg/models"
	"montgomery.wg/http-server/pkg/models/mysql"
)

type logger struct {
	Inner *http.ServeMux
	infoLog *log.Logger
	errorLog *log.Logger
}

type application struct{
	infoLog *log.Logger
	errorLog *log.Logger
	templates *mysql.TemplateModel
}

func (l *logger) ServeHTTP (w http.ResponseWriter, r *http.Request){
	l.infoLog.Printf("Requesting %s?%s\n", r.URL.Path,r.URL.RawQuery)
	l.Inner.ServeHTTP(w,r)
	l.infoLog.Printf("Finished processing request %s\n", r.URL.Path)
}

func (app *application) home (w http.ResponseWriter, r *http.Request) {	
	if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
	if r.Method == "GET" {
		fmt.Fprintf(w, "repo: %s\ntemplate: %s\nversion: %s\n",r.URL.Query().Get("repo"),r.URL.Query().Get("template"),r.URL.Query().Get("version"))
	}else{
		http.NotFound(w, r)
		return
	}
}

func (app *application) showTemplate (w http.ResponseWriter, r *http.Request) {	
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.templates.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    // Write the template data as a plain-text HTTP response body.
    fmt.Fprintf(w, "%v", s)
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/template", app.showTemplate)
	return mux
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := flag.String("dsn", "web:pass@/templates?parseTime=true", "MySQL data source name")
	flag.Parse()
	
	db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
	defer db.Close()
	

	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
		templates: &mysql.TemplateModel{DB: db},
	}

	l := logger{
		Inner: app.routes(),
		infoLog: infoLog,
		errorLog: errorLog,
	}
	
	srv := &http.Server{
		Addr: ":8080",
		ErrorLog: errorLog,
		Handler: &l,
	}

	infoLog.Println("Starting server on port 8080")
	err = srv.ListenAndServe()
	if err != nil{
		log.Panicln(err)
	}
}