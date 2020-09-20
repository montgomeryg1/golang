package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"montgomery.wg/simple-server/pkg/models/books"
)

type application struct {
    errorLog *log.Logger
    infoLog  *log.Logger
    books *books.BooksModel
}

func main() {
	

    addr := ":4000"
    dsn := flag.String("dsn", "web:pass@/books?parseTime=true", "MySQL data source name")
    flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
    //db, err := openDB(*dsn)
    db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
    if err != nil {
        errorLog.Fatal(err)
    }

	app := &application{
        errorLog: errorLog,
        infoLog:  infoLog,
        books: &books.BooksModel{DB: db},
	}
	
    srv := &http.Server{
        Addr:     addr,
        ErrorLog: errorLog,
        Handler:  app.routes(),
    }
    
    // Write messages using the two new loggers, instead of the standard logger.
    infoLog.Printf("Starting server on %s", addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}

