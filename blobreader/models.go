package main

import (
	"time"
)


type BlobName struct {
	Name     string `json:"blobname"`
}

//TruService
type TruService struct {
	TableName       string    `json:"tableName"`
	IPAddress       string    `json:"ipAddress"`
	HTTPMethod      string    `json:"httpMethod"`
	RequestTime     time.Time `json:"requestTime"`
	URL             string    `json:"url"`
	PartnerID       string    `json:"partnerId"`
	MerchantID      string    `json:"merchantId"`
	TerminalID      string    `json:"terminalId"`
	SessionID       string    `json:"sessionId"`
	ResponseStatus  string    `json:"responseStatus"`
	IsCanceled      bool      `json:"isCanceled"`
	ServerTimeMs    float64   `json:"serverTimeMs"`
	ResponseBody    string    `json:"responseBody"`
	ResponseHeaders string    `json:"responseHeaders"`
	RequestBody     string    `json:"requestBody"`
	RequestHeaders  string    `json:"requestHeaders"`
	MachineName     string    `json:"machineName"`
	CorrelationID   string    `json:"correlationId"`
	PartitionKey    string    `json:"partitionKey"`
	RowKey          string    `json:"rowKey"`
	Timestamp       time.Time `json:"timestamp"`
}