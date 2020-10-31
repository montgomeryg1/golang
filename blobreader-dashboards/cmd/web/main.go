package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"montgomery.wg/blobreader/pkg/models/pedapilog"
	"montgomery.wg/blobreader/pkg/models/truservicelog"
)

const (
    host     = "tru-proxy.northeurope.cloudapp.azure.com"
    port     = 6432
    dbname   = "truservicelog"
)



// const (
//     host     = "20.54.91.45"
//     port     = 6432
//     dbname   = "truservicelog"
// )


func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    done := make(chan bool, 1)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    shutdownChannel := make(chan bool, 1)
    ctx := context.Background()
    dbuser := os.Getenv("DB_USER")
    dbpassword := os.Getenv("DB_PASSWORD")

	truserviceDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=truservicelog sslmode=disable", host, port, dbuser, dbpassword)
	truserviceDBpool, err := pgxpool.Connect(ctx, truserviceDSN)
	if err != nil {
		errorLog.Fatalf("Unable to connect to truservicelog database: %v\n", err)
	}
	defer truserviceDBpool.Close()
    truserviceApp := &truservicelog.Application{
        DB: truserviceDBpool,
    }
    conn, err := truserviceDBpool.Acquire(ctx)
    if err = conn.Conn().Ping(ctx); err != nil {
        errorLog.Fatalf("Unable to connect to truservicelog database: %v\n", err)
    }
    defer conn.Conn().Close(ctx)

    pedapiDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=pedapilog sslmode=disable", host, port, dbuser, dbpassword)
    pedapiDBpool, err := pgxpool.Connect(ctx, pedapiDSN)
	if err != nil {
		errorLog.Fatalf("Unable to connect to pedapilog database: %v\n", err)
	}
	defer pedapiDBpool.Close()
    pedapiApp := &pedapilog.Application{
        DB: pedapiDBpool,
	}
    conn, err = pedapiDBpool.Acquire(ctx)
    if err = conn.Conn().Ping(ctx); err != nil {
        errorLog.Fatalf("Unable to connect to peadpi database: %v\n", err)
    }
    defer conn.Conn().Close(ctx)

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthHandler)
    srv := &http.Server{
        Addr:     ":4000",
        ErrorLog: errorLog,
        Handler:  mux,
    }

    go func() {
        infoLog.Printf("Starting server on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            errorLog.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
        }
    }()

    
    
    var count int = 1
    go func() {
        
        var wg sync.WaitGroup
        for {
            
            select {
            case <-shutdownChannel:
                infoLog.Println("Dashboards shutdown...")
                close(shutdownChannel)
                return
    
            case <-time.After(time.Second*60):
                infoLog.Println("Loop:", count)
                count++
                break
            }
            
            hour, minute, _ := time.Now().UTC().Clock()
            
            wg.Add(1)
            go func() {
                defer wg.Done()
                if minute == 10 {
                    infoLog.Println("Updating ratings tables...")
                    if err := truserviceApp.GetRatingsByHour(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished updating ratings tables...")
                }
            }()
            
            wg.Add(1)
            go func() {
                defer wg.Done()
                if minute == 30 {
                    infoLog.Println("Updating world map table...")
                    if err := truserviceApp.WorldMap(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished updating world map table for truservicelog...")

                    infoLog.Println("Starting to update world map table for pedapilog...")
                    if err := pedapiApp.WorldMap(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished updating world map table for pedapilog...")
                }
            }()

            wg.Add(1)
            go func() {
                defer wg.Done()
                if hour == 6 && minute == 15 {
                    infoLog.Println("Updating lastseen table...")
                    if err := truserviceApp.MerchantLastSeen(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished updating lastseen table...")
                }
            }()

            
            wg.Add(1)
            go func() {
                defer wg.Done()
                if (minute % 2) == 0 {
                    infoLog.Println("Deleting old records...")
                    _, err := truserviceApp.DB.Exec(ctx, "DELETE FROM httplogs where requesttime < NOW() - INTERVAL '8 day'")
                    if err != nil {
                        errorLog.Println(err)
                    }
                    _, err = truserviceApp.DB.Exec(ctx, "DELETE FROM worldmap where requesttime < NOW() - INTERVAL '3 day'")
                    if err != nil {
                        errorLog.Println(err)
                    }
                    _, err = pedapiApp.DB.Exec(ctx, "DELETE FROM httplogs where requesttime < NOW() - INTERVAL '8 day'")
                    if err != nil {
                        errorLog.Println(err)
                    }
                    _, err = pedapiApp.DB.Exec(ctx, "DELETE FROM worldmap where requesttime < NOW() - INTERVAL '3 day'")
                    if err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished deleting old records...")
                }
            }()


            wg.Add(1)
            go func() {
                defer wg.Done()
                if hour == 1 && minute == 15 {
                    infoLog.Println("Updating requests by partnerid table...")
                    if err := truserviceApp.RequestsByPartnerID(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    if err := pedapiApp.RequestsByPartnerID(ctx); err != nil {
                        errorLog.Println(err)
                    }
                    infoLog.Println("Finished requests by partnerid table...")
                }
            }()

            wg.Wait()            
            // deadline,_:=cancelCtx.Deadline()
            // truserviceApp.InfoLog.Println("Time until context timeout:", time.Until(deadline))
            
        }
    }()

    go waitForShutdown(srv, errorLog, quit, done, shutdownChannel)

    <-done
    infoLog.Println("Server stopped")
    os.Exit(0)
}



func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello from Blobreader Dashboards!"))
}

func waitForShutdown(srv *http.Server, errorLog *log.Logger, quit <-chan os.Signal, done chan<- bool, shutdownChannel chan bool) {
    <-quit
    errorLog.Println("Server is shutting down...")
    shutdownChannel <- true
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
  
    srv.SetKeepAlivesEnabled(false)
    if err := srv.Shutdown(ctx); err != nil {
      errorLog.Fatalf("Could not gracefully shutdown the server: %v\n", err)
    }
    close(done)
}