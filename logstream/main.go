package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

func stream(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    cmd := exec.Command("curl","-u","$tru-live-profiles-api-northeurope:iKo6NzfdDCdktqGHKLHmW5SLEcqTo2eatbjQjwiNMAQ4qTFegDbbauYe8PkP","https://tru-live-profiles-api-northeurope.scm.azurewebsites.net/api/logstream")
    rPipe, wPipe, err := os.Pipe()
    if err != nil {
        log.Fatal(err)
    }
    cmd.Stdout = wPipe
    cmd.Stderr = wPipe
    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }
    go writeOutput(w, rPipe)
    cmd.Wait()
    wPipe.Close()
}

func writeOutput(w http.ResponseWriter, input io.ReadCloser) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }

    // Important to make it work in browsers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    in := bufio.NewScanner(input)
    for in.Scan() {
        fmt.Fprintf(w, "%s\n\n", in.Text())
        flusher.Flush()
    }
    input.Close()
}

func main() {
	


	router := httprouter.New()
    router.GET("/", stream)

    if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}