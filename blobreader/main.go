package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger
	db *pgxpool.Pool
	containerURL *azblob.ContainerURL
}

func ( app *application) CheckError(err error) {
    if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
        app.errorLog.Output(2, trace)
    }
}


func ( a *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	// if err := a.db.Ping(); err != nil {
    //     w.WriteHeader(http.StatusInternalServerError)
    // }
	w.WriteHeader(http.StatusOK)
}

func waitForShutdown(srv *http.Server,chSrch chan string, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	
	close(chSrch)
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	srv.Shutdown(ctx)

	close(done)
}

func main() {

	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	storagekey := os.Getenv("STORAGE_KEY")
	storageAcct := os.Getenv("STORAGE_ACCOUNT")
	blobContainer := os.Getenv("BLOB_CONTAINER")

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	blobs := make(chan string, 25)
	processedBlobs := make(chan string)
	//var wg sync.WaitGroup

	ctx := context.Background()

	chSrch := make(chan string)
	go timeString(chSrch)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, dbuser, dbpassword, dbname)
	dbpool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	credential, err := azblob.NewSharedKeyCredential(storageAcct, storagekey)
	if err != nil {
        errorLog.Fatal(err)
    }
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", storageAcct))
	serviceURL := azblob.NewServiceURL(*u, p)
	containerURL := serviceURL.NewContainerURL(blobContainer)

	app := &application{
        errorLog: errorLog,
		infoLog:  infoLog,
		db: dbpool,
		containerURL: &containerURL,
	}

	//Create worker pool
	for i := 0; i < cap(blobs); i++ {
		//go worker(ctx, blobs, &wg, app)
		go worker(ctx, blobs, processedBlobs, app)
	}


	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", app.healthHandler)
	mux.Handle("/metrics", promhttp.Handler())


	srv := &http.Server{
		Handler:      mux,
		Addr:         ":4000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

    go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			errorLog.Fatal(err)
		}
	}()

	go waitForShutdown(srv, chSrch, quit, done)
	
	var blobcount int
	for n := range chSrch {
		for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
			//fmt.Println("marker: ", marker)
			//app.infoLog.Println("searchstring: ", n)
			// Get a result segment starting with the blob indicated by the current Marker.
			listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{Prefix: n })
			app.CheckError(err)
			// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
			// the next segment (after processing the current result segment).
			marker = listBlob.NextMarker

			// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
			blobcount = len(listBlob.Segment.BlobItems)
			infoLog.Printf("Processing %d blobs", len(listBlob.Segment.BlobItems))
			go func (){
				for _, blobInfo := range listBlob.Segment.BlobItems {
					//wg.Add(1)
					blobs <- blobInfo.Name
				}
			}()

		}
		for i:=0; i < blobcount; i++ {
			<- processedBlobs
			//infoLog.Println("Number of blobs processed:", i + 1)
		}
		infoLog.Printf("Finished processing %d blobs", blobcount)

	}
	
	close(blobs)
	<- done
	infoLog.Println("Server shutting down")
}