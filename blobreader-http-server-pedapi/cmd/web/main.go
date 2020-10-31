package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"montgomery.wg/pedapilog/pkg/models/pgsql"
)


const (
    host     = "tru-proxy.northeurope.cloudapp.azure.com"
    port     = 6432
    user     = "pgsqladmin@tru-dev-main-pgsql-eun001"
    dbname   = "pedapilog"
)

type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger
    httprequests *pgsql.PEDAPIRecordModel
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    done := make(chan bool, 1)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    
	dbpassword := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, dbpassword, dbname)
	db, err := openDB(dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    app := &application{
        errorLog: errorLog,
        infoLog:  infoLog,
        httprequests: &pgsql.PEDAPIRecordModel{DB: db},
	}
	
    srv := &http.Server{
        Addr:     ":8080",
        ErrorLog: errorLog,
        Handler:  app.routes(),
    }

    go waitForShutdown(srv, errorLog, quit, done)

    infoLog.Printf("Starting server on %s", srv.Addr)

    if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        errorLog.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
    }

    <-done
    infoLog.Println("Server stopped")
}


func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

func waitForShutdown(srv *http.Server, errorLog *log.Logger, quit <-chan os.Signal, done chan<- bool) {
    <-quit
    errorLog.Println("Server is shutting down...")
  
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
  
    srv.SetKeepAlivesEnabled(false)
    if err := srv.Shutdown(ctx); err != nil {
      errorLog.Fatalf("Could not gracefully shutdown the server: %v\n", err)
    }
    close(done)
}