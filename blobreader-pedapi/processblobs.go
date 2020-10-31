package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func worker(ctx context.Context, blobs chan string, processedBlobs chan string, app *application) {
	insertDynStmt := `insert into "httplogs"(requestheaders, requestbody, responseheaders, responsebody, requesttime, httpmethod, url, responsestatus, servertimems, ipaddress, partnerid, merchantid, terminalid, blobname) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ON CONFLICT ON CONSTRAINT blob_unique DO NOTHING`
	for b := range blobs {
		//app.infoLog.Println("Blob name: ", b)

		blobURL := app.containerURL.NewBlockBlobURL(b)
		get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		downloadedData := &bytes.Buffer{}
		reader := get.Body(azblob.RetryReaderOptions{})
		downloadedData.ReadFrom(reader)
		reader.Close() // The client must close the response body when finished with it
		data := PEDAPI{}

		_ = json.Unmarshal(downloadedData.Bytes(), &data)

		_, err = app.db.Exec(ctx,insertDynStmt, data.RequestHeaders, data.RequestBody, data.ResponseHeaders, data.ResponseBody, data.RequestTime, data.HTTPMethod, data.URL, data.ResponseStatus, data.ServerTimeMs, data.IPAddress, data.PartnerID, data.MerchantID, data.TerminalID, b)
		if err != nil {
			app.errorLog.Println(err)
		}

		//wg.Done()
		processedBlobs <- b
	}
}