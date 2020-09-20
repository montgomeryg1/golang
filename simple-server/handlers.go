package main

import (
	"encoding/json"
	"net/http"

	"montgomery.wg/simple-server/pkg/models"
)

func (app *application) check(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("GET params were:", r.URL.Query())
	var params = make(map[string]int)
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var total = 0
	
	for i := 1;  i<=10; i++ {
		resp, err := http.Get(id)
		if err != nil {
			app.errorLog.Println(err.Error())
			params[id] = http.StatusNotFound
		}else{
			app.infoLog.Printf("Returned status code %d", resp.StatusCode)
			total = total + resp.StatusCode
			params[id] = resp.StatusCode
		}
	}


	app.infoLog.Printf("Total: %d", total)
	if total == 2000 {
		w.WriteHeader(http.StatusOK)
	}else{
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(&params)

}


func (app *application) addBook(w http.ResponseWriter, r *http.Request) {
    // if r.Method != http.MethodPost {
    //     w.Header().Set("Allow", http.MethodPost)
    //     app.clientError(w, http.StatusMethodNotAllowed)
    //     return
    // }

    title := "Start with why"
    content := "Simon Sinek"
	book := models.Book{Title: title, Author: content}
    app.books.DB.Create(&book)
}