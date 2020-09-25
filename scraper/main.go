package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var (
	fileName    string
	fullUrlFile string
)

func main() {

	// fmt.Println(os.Args[1])
	// urlName := os.Args[1]
	scrape()

	// fullUrlFile = "https://unsplash.com/photos/lHJZ2JTxSlI/download?force=true"

	// Build fileName from fullPath
	// buildFileName()

	// Create blank file
	// file := createFile()

	// Put content on file
	// putFile(file, httpClient())

}

func putFile(file *os.File, client *http.Client) {
	resp, err := client.Get(fullUrlFile)

	checkError(err)

	defer resp.Body.Close()

	//fmt.Println(resp.Header)
	// fmt.Println(resp.Header.Get("Content-Disposition"))
	// _, params, err := mime.ParseMediaType(`resp.Header.Get("Content-Disposition")`)
	// fmt.Println(resp.Header.Get("Content-Disposition"))
	// filename := params["filename"]
	if strings.Index(resp.Header.Get("Content-Type"), "image/jpeg") > -1 {
		size, err := io.Copy(file, resp.Body)

		defer file.Close()

		checkError(err)

		fmt.Printf("Just Downloaded a file %s with size %d \n", fileName, size)
	}
}

func buildFileName() {
	fileUrl, err := url.Parse(fullUrlFile)
	checkError(err)

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName = "./images/" + segments[len(segments)-2] + ".jpg"
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFile() *os.File {
	file, err := os.Create(fileName)

	checkError(err)
	return file
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func scrape() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("unsplash.com"),
		// Visit only root url and urls which start with "e" or "h" on httpbin.org
		// colly.URLFilters(
		// 	// regexp.MustCompile("http://httpbin\\.org/(|z.+)$"),
		// 	// regexp.MustCompile("http://httpbin\\.org/h.+"),
		// 	// regexp.MustCompile(".*photos.*"),
		// 	// regexp.MustCompile("(download?force=true)"),
		// 	// regexp.MustCompile("https?://unsplash\\.com/t/wallpapers"),
		// 	regexp.MustCompile("https?://unsplash\\.com/(|photos/.+)$"),
		// // regexp.MustCompile(u),
		// // regexp.MustCompile("https?://unsplash\\.com/photos/(.+)force=true$"),
		// // regexp.MustCompile("https?://.+$"),
		// ),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Contains(link, "download") {
			// fmt.Println(e.Request.AbsoluteURL(link))
			c.Visit(e.Request.AbsoluteURL(link))

		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		fullUrlFile = r.URL.String()
		buildFileName()
		file := createFile()
		putFile(file, httpClient())
	})

	c.Visit("https://unsplash.com/t/wallpapers")
}
