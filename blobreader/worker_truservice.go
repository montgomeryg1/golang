package main

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/jackc/pgx/v4"
)

func worker_truservice(ctx context.Context, blobs chan string, processedBlobs chan string, dataBase *db, containerURL *azblob.ContainerURL, batch *pgx.Batch) {
	receiptRex := regexp.MustCompile(`(?is)<Receipt>(.*?)</Receipt>`)
	messageRex := regexp.MustCompile(`(?is)<Message>(.*?)</Message>`)
	questionRex := regexp.MustCompile(`(?is)<Question TimeoutMs="\d+">(.*?)(4|0)`)
	ratingRex := regexp.MustCompile(`Rating Value="(.*?)"`)
	ipRex := regexp.MustCompile(`X-Azure-ClientIP: (25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	//insertDynStmt := `insert into "httplogs"(blobname, requesttime, messagetype, infomessage, sessionid, ipaddress, httpmethod, responsestatus, partnerid, merchantid, terminalid, url, requestbody, requestheaders, responsebody, responseheaders) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) ON CONFLICT ON CONSTRAINT blob_unique	DO NOTHING`
	for b := range blobs {
		//app.infoLog.Println("Blob name: ", b)

		blobURL := containerURL.NewBlockBlobURL(b)
		get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
		dataBase.CheckError(err)
		downloadedData := &bytes.Buffer{}
		reader := get.Body(azblob.RetryReaderOptions{})
		downloadedData.ReadFrom(reader)
		reader.Close() // The client must close the response body when finished with it
		data := TruService{}
		_ = json.Unmarshal(downloadedData.Bytes(), &data)

		var infoMsg string = "."
		var msgType string
		var str string
		frmtReq := strings.ReplaceAll(data.RequestBody, "><", ">\n<")
		frmtRes := strings.ReplaceAll(data.ResponseBody, "><", ">\n<")
		ipAddr := data.IPAddress
		submatchall := ipRex.FindAllString(data.RequestHeaders, -1)
		if submatchall != nil {
			for _, element := range submatchall {
				ipAddr = strings.TrimPrefix(element, "X-Azure-ClientIP: ")
			}
		}
		switch {
		case strings.Contains(data.RequestBody, "</PosEventList>"):
			//fmt.Println("Writing PosEventList record")
			msgType = "PosEventList"
		case strings.Contains(data.RequestBody, "</Query>"):
			//fmt.Println("Writing Query record")
			msgType = "Query"
			if data.ResponseStatus == "200" {
				out := receiptRex.FindAllStringSubmatch(data.ResponseBody, -1)
				for _, i := range out {
					str += i[1]
				}
				infoMsg = strings.TrimSpace(strings.ReplaceAll(str, "\\n", " "))
			} else {
				out := messageRex.FindAllStringSubmatch(data.ResponseBody, -1)
				for _, i := range out {
					str += i[1]
				}
				infoMsg = strings.TrimSpace(strings.ReplaceAll(str, "\\n", " "))
			}
		case strings.Contains(data.RequestBody, "</Question>"):
			//fmt.Println("Writing Question record")
			msgType = "Question"
			out := questionRex.FindAllStringSubmatch(data.ResponseBody, -1)
			for _, i := range out {
				str += i[1]
			}
			infoMsg = strings.ReplaceAll(str, "\\n", " ")
			if data.ResponseStatus == "401" || data.ResponseStatus == "400" {
				out := messageRex.FindAllStringSubmatch(data.ResponseBody, -1)
				for _, i := range out {
					str += i[1]
				}
				infoMsg = strings.TrimSpace(strings.ReplaceAll(str, "\\n", " "))
			}
		case strings.Contains(data.RequestBody, "</Rating>"):
			//fmt.Println("Writing Rating record")
			msgType = "Rating"
			out := ratingRex.FindAllStringSubmatch(data.RequestBody, -1)
			for _, i := range out {
				str += i[1]
			}
			//fmt.Println(strings.TrimSpace(infoMsg))
			if data.ResponseStatus == "401" || data.ResponseStatus == "400" {
				str = ""
				out := messageRex.FindAllStringSubmatch(data.ResponseBody, -1)
				for _, i := range out {
					str += i[1]
				}
			}
			infoMsg = strings.ReplaceAll(str, "\\n", " ")
		case strings.Contains(data.RequestBody, "<Transaction"):
			//fmt.Println("Writing Transaction record")
			msgType = "Transaction"
		case strings.Contains(data.RequestBody, "<Rating DateTime="):
			//fmt.Println("Writing Rating record")
			msgType = "Rating"
		default:
			msgType = "Unknown"
		}
		batch.Queue("insert into httplogs(blobname, requesttime, messagetype, infomessage, sessionid, ipaddress, httpmethod, responsestatus, partnerid, merchantid, terminalid, url, requestbody, requestheaders, responsebody, responseheaders) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) ON CONFLICT ON CONSTRAINT blob_unique DO NOTHING", b, data.RequestTime, msgType, infoMsg, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
		//_, e := app.db.Exec(ctx, insertDynStmt, b, data.RequestTime,  msgType, infoMsg, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
		//app.CheckError(e)
		processedBlobs <- b
	}

}
