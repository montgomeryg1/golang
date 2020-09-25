package main

import (
  "database/sql"
  "fmt"

  _ "github.com/lib/pq"
)

const (
  host     = "tru-dev-main-pgsql-eun001.postgres.database.azure.com"
  port     = 5432
  user     = "sonarqube@tru-dev-main-pgsql-eun001"
  password = "T5haaGYkPg196E7StkovW5U5"
  dbname   = "sonarDB"
)

func main() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")
}