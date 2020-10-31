package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type PEDAPIRecord struct {
	RequestTime     time.Time `json:"requestTime"`
	URL             string    `json:"url"`
	PartnerID       string    `json:"partnerId"`
	MerchantID      string    `json:"merchantId"`
	TerminalID      string    `json:"terminalId"`
	ResponseStatus  string    `json:"responseStatus"`
	ResponseBody    string    `json:"responseBody"`
	RequestBody     string    `json:"requestBody"`
}
