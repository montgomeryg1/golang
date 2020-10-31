package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"montgomery.wg/blobreader/pkg/models/truservicelog"
)


const (
    host     = "tru-proxy.northeurope.cloudapp.azure.com"
    port     = 6432
    user     = "pgsqladmin@tru-dev-main-pgsql-eun001"
    dbname   = "truservicelog"
)


func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    done := make(chan bool, 1)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    ctx := context.Background()
    dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	truserviceDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=truservicelog sslmode=disable", host, port, dbuser, dbpassword)
	truserviceDBpool, err := pgxpool.Connect(ctx, truserviceDSN)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer truserviceDBpool.Close()


    app := &truservicelog.Application{
        errorLog: errorLog,
        infoLog:  infoLog,
        DB: truserviceDBpool,
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