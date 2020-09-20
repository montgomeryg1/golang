package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	port := getPort()
	r.HandleFunc("/", indexHandler).Methods("GET")
	err := http.ListenAndServe(port, r)

	if err != nil {
		log.Fatal("Error listening on server:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Got my server up and running in Go. Yay!!!!")

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3500"
		fmt.Printf("PORT NOT DEFINED. USING THE PORT %s as the running port\n", port)
	}
	return port
}
