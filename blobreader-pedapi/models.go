package main

import (
	"time"
)


type BlobName struct {
	Name     string `json:"blobname"`
}

type PEDAPI struct {
	TableName       string    `json:"tableName"`
	RequestHeaders  string    `json:"requestHeaders"`
	RequestBody     string `json:"requestBody"`
	ResponseHeaders string    `json:"responseHeaders"`
	ResponseBody    string    `json:"responseBody"`
	RequestTime     time.Time `json:"requestTime"`
	HTTPMethod      string    `json:"httpMethod"`
	URL             string    `json:"url"`
	ResponseStatus  string    `json:"responseStatus"`
	ServerTimeMs    float64   `json:"serverTimeMs"`
	IsCanceled      bool      `json:"isCanceled"`
	IPAddress       string    `json:"ipAddress"`
	PartnerID       string    `json:"partnerId"`
	MerchantID      string    `json:"merchantId"`
	TerminalID      string    `json:"terminalId"`
	ClientID        int       `json:"clientId"`
	PartitionKey    string    `json:"partitionKey"`
	RowKey          string    `json:"rowKey"`
	Timestamp       time.Time `json:"timestamp"`
}