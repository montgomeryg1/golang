package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}


// func ConnectDataBase() {
// 	database, err := gorm.Open("sqlite3", "test.db")

// 	if err != nil {
// 		panic("Failed to connect to database!")
// 	}

// 	database.AutoMigrate(&Book{})

// 	DB = database
// }