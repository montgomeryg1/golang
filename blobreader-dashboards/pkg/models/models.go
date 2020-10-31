package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type TruserviceRequest struct {
	TableName       string    `json:"tableName"`
	IPAddress       string    `json:"ipAddress"`
	HTTPMethod      string    `json:"httpMethod"`
	RequestTime     time.Time `json:"requestTime"`
	URL             string    `json:"url"`
	PartnerID       string    `json:"partnerId"`
	MerchantID      string    `json:"merchantId"`
	TerminalID      string    `json:"terminalId"`
	SessionID       string    `json:"sessionId"`
	MessageType		string	  `json:"messagetype"`
	InfoMessage		string	  `json:"infomessage"`
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


// func NewPedapiPool(ctx context.Context, pedapiDSN string, errorLog *log.Logger, infoLog *log.Logger) (*pedapilog.Application) {
// 	pedapiDBpool, err := pgxpool.Connect(ctx, pedapiDSN)

// 	if err != nil {
// 		errorLog.Fatalf("Unable to connect to database: %v\n", err)
// 	}

// 	pedapiApp := &pedapilog.Application{
// 		ErrorLog: errorLog,
// 		InfoLog:  infoLog,
// 		DB:       pedapiDBpool,
// 	}
// 	return pedapiApp
// }

// func NewTruServicePool(ctx context.Context, truserviceDSN string, errorLog *log.Logger, infoLog *log.Logger) (*truservicelog.Application) {
// 	truserviceDBpool, err := pgxpool.Connect(ctx, truserviceDSN)
// 	if err != nil {
// 		errorLog.Fatalf("Unable to connect to database: %v\n", err)
// 	}

// 	truserviceApp := &truservicelog.Application{
// 		ErrorLog: errorLog,
// 		InfoLog:  infoLog,
// 		DB:       truserviceDBpool,
// 	}
// 	return truserviceApp
// }