package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Template struct {
    ID          int
    Repo        string
    Template    string
    Version     string
    Created     time.Time
    Expires     time.Time
}