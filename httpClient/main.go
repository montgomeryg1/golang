package main

import (
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    url := "https://api.globalcode.com.br/v1/publico/eventos"

    // Create a Bearer string by appending string access token
    var bearer = "Bearer " + <ACCESS TOKEN HERE>

    // Create a new request using http
    req, err := http.NewRequest("GET", url, nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERRO] -", err)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    log.Println(string([]byte(body)))
}