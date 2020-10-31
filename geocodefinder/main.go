//

// csvFile, _ := os.Open("region.csv")
// r := csv.NewReader(csvFile)
// for {
//   record, err := r.Read()
//   if err == io.EOF {
// 	break
//   }
//   if err != nil {
// 	log.Fatal(err)
//   }
//   //fmt.Println(record)
//   //query := "Yorkville Ohio United States"
//   query := strings.Join(record," ")
//   lat, lng, err := geocoder.Geocode(query)
//   if err != nil {
// 	panic("THERE WAS SOME ERROR!!!!!")
//   }
//   fmt.Printf("\"%v\",\"%f\",\"%f\",\n", query, lat, lng)
// }
//}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	url = "https://freegeoip.app/json"
	host     = "tru-proxy.northeurope.cloudapp.azure.com"
    port     = 6432
    user     = "pgsqladmin@tru-dev-main-pgsql-eun001"
)

type Client struct {}

type IPloc struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

type GeoRecord struct {
	IP          string  `json:"ip"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	City        string  `json:"city"`
	CountryName string  `json:"country_name"`
}


type IPRecord struct {
	IP          string
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello from Blobreader-dashboards"))
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

func (a *application) CheckError(err error) {
    if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    	a.errorLog.Output(2, trace)
        a.errorLog.Fatal(err)
    }
}

func (c *Client) GeoInfo(q string) (*IPloc, error){
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s",url, q), nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret IPloc
	if err := json.NewDecoder(res.Body).Decode(&ret);err !=nil{
		return nil, err
	}
	return &ret, nil
}

type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger
}


func main() {

	
	dbpassword := os.Getenv("DB_PASSWORD")
	//dbname := os.Getenv("DB_NAME")

	ctx := context.Background()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	selectDynStmt := `select ipaddress from "geoloc" where ipaddress = $1`
	insertDynStmt := `insert into "geoloc"(ipaddress, latitude, longitude, city) values($1, $2, $3,  $4)`
	distinctIPStmt := `SELECT distinct ipaddress from httplogs WHERE httplogs.requesttime >= NOW() - '1 hour'::INTERVAL`

	done := make(chan bool, 1)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthHandler)
	
	app := &application{
        errorLog: errorLog,
        infoLog:  infoLog,
	}
	
	var c Client

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

    go waitForShutdown(srv, errorLog, quit, done)


	
	go func() {
		for {
			dataBases := []string{"truservicelog", "pedapilog"}
			for i, dbname := range dataBases {
				fmt.Println(i, dbname)
		
				// connection string
				dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, dbpassword, dbname)
				dbpool, err := pgxpool.Connect(ctx, dsn)
				if err != nil {
					errorLog.Fatalf("Unable to connect to database: %v\n", err)
				}
		
				rows, err := dbpool.Query(ctx, distinctIPStmt)
				app.CheckError(err)
			
				for rows.Next() {
					//distinctIP := &IPRecord{}
					var distinctIP pgtype.Inet
					var checkIP pgtype.Inet
					err = rows.Scan(&distinctIP)
					app.CheckError(err)
					//var record GeoRecord
					err := dbpool.QueryRow(ctx, selectDynStmt, &distinctIP.IPNet.IP).Scan(&checkIP)
					switch {
						case err == pgx.ErrNoRows:
							app.infoLog.Printf("%s is being added to table", distinctIP.IPNet.IP.String())
							info, _ := c.GeoInfo(distinctIP.IPNet.IP.String())
							// app.CheckError(err)
							if info.Latitude != 0.000000 && info.Longitude != 0.000000 {
								if info.City != ""{
									_, err := dbpool.Exec(ctx, insertDynStmt, &distinctIP.IPNet.IP, info.Latitude, info.Longitude, info.City)
									app.CheckError(err)
								}else{
									_, err := dbpool.Exec(ctx, insertDynStmt, &distinctIP.IPNet.IP, info.Latitude, info.Longitude, info.CountryName)
									app.CheckError(err)
								}
								
							}					
						case err != nil:
							app.errorLog.Fatal(err)
						default:
							app.infoLog.Printf("%s already in table\n", distinctIP.IPNet.IP.String())
					}
			
				}
		
				rows.Close()
				dbpool.Close()
			}
			time.Sleep(30 * time.Minute)
		}
	}()

	<-done
    infoLog.Println("Server stopped")
}
