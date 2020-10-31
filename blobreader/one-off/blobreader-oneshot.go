package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	_ "github.com/lib/pq"
)

const (
    port     = 6432
)


// const (
//     host     = "20.54.91.45"
//     port     = 6432
//     user     = "pgbench"
//     dbname   = "truservicelog"
// 	dbpassword = "5g^AXgB6nt^9SOX1H7ixuvOF"
// )

type Truservicejson struct {
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

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, dbpassword, dbname)
         
	// open database
    db, err := sql.Open("postgres", psqlconn)
    CheckError(err)
     
	// close database
    defer db.Close()
 
	// check db
    err = db.Ping()
    CheckError(err)
 
	infoLog.Println("Connected to database!")	

	//insertStmt := `insert into "httplogs"("RequestTime", "PartnerID", "MerchantID", "TerminalID", "RequestBody", "RequestHeaders", "ResponseBody", "ResponseHeaders") values('2020-10-05T00:01:59.912127Z','4300','1103','3','<?xml version=\"1.0\"?>\n<Request PartnerId=\"4300\" MerchantId=\"1103\" TerminalId=\"3\" SessionId=\"129919\" xmlns=\"http://docs.trurating.com/schema/truservice/v220.xsd\"><Rating Value=\"4\" ResponseTimeMs=\"6102\" Rfc1766=\"en\" DateTime=\"2020-10-04T19:01:44\"><Transaction Type=\"SALE\" Id=\"480300038\" DateTime=\"2020-10-04T19:01:44\" Amount=\"1407\" Gratuity=\"000\" Currency=\"840\" Result=\"APPROVED\"><Tender CardType=\"MASTERCARD\" EntryMode=\"05\" TenderType=\"CREDIT\" Amount=\"1407\"><CardHash Type=\"TRUTRACE\" Value=\"A0000000041010#MASTERCARD#5089#191210###01#185592E533AF3678D533EA26075C9249D984BB61E344F3B647BA842723A0E7BB#2996382021299BEB5FC27A4714AD90FF634DE48B372279DA9FFC01C5973889E7\"/></Tender></Transaction></Rating></Request>\n','Connection: Keep-Alive\r\nVia: 1.1 Azure\r\nAccept: */*\r\nHost: service-v2xx-southcentralus.trurating.com\r\nMax-Forwards: 10\r\nx-tru-api-partner-id: 4300\r\nx-tru-api-merchant-id: 1103\r\nx-tru-api-terminal-id: 3\r\nx-tru-api-encryption-scheme: 3\r\nx-tru-api-mac: 583FE99D7DBA1DBDD916BCB0D532B9CC2580C39C7501FAF5AB99892B2D457705DFFE88532C99577D\r\nX-Forwarded-For: 108.208.92.145, 147.243.144.166:33258\r\nX-Azure-ClientIP: 108.208.92.145\r\nX-Azure-Ref: 0d2J6XwAAAAAA8KmsSLr1RINB5rG0bEetREFMRURHRTEwMDUAY2U4ODQ2MzMtMzcyNi00YWNmLWE0MTktY2UxMjhlYTNkZDgy\r\nX-Forwarded-Host: service-v2xx.trurating.com\r\nX-Forwarded-Proto: https\r\nX-Azure-RequestChain: hops=1\r\nX-Azure-SocketIP: 108.208.92.145\r\nX-Azure-FDID: ce884633-3726-4acf-a419-ce128ea3dd82\r\nX-WAWS-Unencoded-URL: /api/servicemessage\r\nCLIENT-IP: 147.243.144.166:33258\r\nX-ARR-LOG-ID: 948bda18-f829-475b-836a-84d441fef908\r\nDISGUISED-HOST: service-v2xx-southcentralus.trurating.com\r\nX-SITE-DEPLOYMENT-ID: tru-live-service-v2xx-southcentralus__8af0\r\nWAS-DEFAULT-HOSTNAME: tru-live-service-v2xx-southcentralus.azurewebsites.net\r\nX-Original-URL: /api/servicemessage\r\nX-ARR-SSL: 2048|256|C=US, O=DigiCert Inc, OU=www.digicert.com, CN=Thawte RSA CA 2018|C=GB, L=London, O=TruRating Limited, CN=*.trurating.com\r\nX-AppService-Proto: https\r\nX-Forwarded-TlsVersion: 1.2\r\n','<?xml version=\"1.0\"?>\r\n<Response PartnerId=\"4300\" MerchantId=\"1103\" TerminalId=\"3\" SessionId=\"129919\" xmlns=\"http://docs.trurating.com/schema/truservice/v220.xsd\" />',' ')`
    //_, e := db.Exec(insertStmt)
	insertDynStmt := `insert into "httplogs"(blobname, requesttime, messagetype, infomessage, sessionid, ipaddress, httpmethod, responsestatus, partnerid, merchantid, terminalid, url, requestbody, requestheaders, responsebody, responseheaders) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	storagekey := os.Getenv("STORAGE_KEY")
	storageAcct := os.Getenv("STORAGE_ACCOUNT")
	credential, err := azblob.NewSharedKeyCredential(storageAcct, storagekey)
	if err != nil {
		errorLog.Fatal(err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", "trulivelogsneurope"))

	serviceURL := azblob.NewServiceURL(*u, p)

	ctx := context.Background()

	containerURL := serviceURL.NewContainerURL(dbname)

	msgSt := regexp.MustCompile(`<Message>`)
	msgEnd := regexp.MustCompile(`</Message>`)
	ip := regexp.MustCompile(`X-Azure-ClientIP: (25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	for i := 8; i < 18; i++ {
		n := fmt.Sprintf("2020/10/17/%s", strconv.Itoa(i))
		infoLog.Println(n)
		for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
			//fmt.Println("marker: ", marker)
			// Get a result segment starting with the blob indicated by the current Marker.
			listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{Prefix: n})
			if err != nil {
				errorLog.Fatal(err)
			}
			// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
			// the next segment (after processing the current result segment).
			marker = listBlob.NextMarker


			// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
			
			for _, blobInfo := range listBlob.Segment.BlobItems {
					infoLog.Println("Blob name: ", blobInfo.Name)
					blobURL := containerURL.NewBlockBlobURL(blobInfo.Name)
					get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
					if err != nil {
						errorLog.Fatal(err)
					}
					downloadedData := &bytes.Buffer{}
					reader := get.Body(azblob.RetryReaderOptions{})
					downloadedData.ReadFrom(reader)
					reader.Close() // The client must close the response body when finished with it
					data := Truservicejson{}
					_ = json.Unmarshal(downloadedData.Bytes(), &data)
					qryRes := "no message"
					frmtReq := strings.ReplaceAll(data.RequestBody, "><",">\n<")
					frmtRes := strings.ReplaceAll(data.ResponseBody, "><",">\n<")
					ipAddr := data.IPAddress
					submatchall := ip.FindAllString(data.RequestHeaders, -1)
					if submatchall != nil {
						for _, element := range submatchall {
							ipAddr = strings.TrimPrefix(element,"X-Azure-ClientIP: ")
						}					
					}

					switch {
						case strings.Contains(data.RequestBody,"</PosEventList>"):
							//fmt.Println("Writing PosEventList record")
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime,  "PosEventList", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"</Query>"):
							//fmt.Println("Writing Query record")
							
							if data.ResponseStatus == "200"{
								switch{
									case strings.Contains(data.ResponseBody,"SUSPEND"):
										qryRes = "TruHost response is SUSPEND"
									case strings.Contains(data.ResponseBody,"DEACTIVATE"):
										qryRes = "TruHost response is DEACTIVATE"
									default:
										qryRes = "TruHost response is OK"
								}
							}else{
								result := msgSt.Split(data.ResponseBody, 2)
								finalresult := msgEnd.Split(result[1],2)
								qryRes = finalresult[0]
							}
							
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Query", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"</Question>"):
							//fmt.Println("Writing Question record")
							if data.ResponseStatus == "401"{
								result := msgSt.Split(data.ResponseBody, 2)
								finalresult := msgEnd.Split(result[1],2)
								qryRes = finalresult[0]
							}						
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Question", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)							
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"</Rating>"):
							//fmt.Println("Writing Rating record")
							if data.ResponseStatus == "400"{
								result := msgSt.Split(data.ResponseBody, 2)
								finalresult := msgEnd.Split(result[1],2)
								qryRes = finalresult[0]
							}						
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Rating", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"</Transaction>"):
							//fmt.Println("Writing Transaction record")
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Transaction", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)							
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"<Transaction Type="):
							//fmt.Println("Writing Transaction record")
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Transaction", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)							
							if e != nil {
								errorLog.Println(e)
							}
						case strings.Contains(data.RequestBody,"<Rating DateTime="):
							//fmt.Println("Writing Rating record")
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Rating", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)							
							if e != nil {
								errorLog.Println(e)
							}
						default:
							_, e := db.Exec(insertDynStmt, blobInfo.Name, data.RequestTime, "Unknown", qryRes, data.SessionID, ipAddr, data.HTTPMethod, data.ResponseStatus, data.PartnerID, data.MerchantID, data.TerminalID, data.URL, frmtReq, data.RequestHeaders, frmtRes, data.ResponseHeaders)
							if e != nil {
								errorLog.Println(e)
							}
					}
			}
		}
	}
	
}