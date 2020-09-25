package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Offers struct {
	Type  string `json:"@type"`
	Price string `json:"price"`
}

type WorkExample struct {
	Type                 string `json:"@type"`
	Name                 string `json:"name"`
	DatePublished        string `json:"datePublished"`
	Description          string `json:"description"`
	Duration             string `json:"duration"`
	RequiresSubscription string `json:"requiresSubscription"`
}

type Podcast struct {
	Context       string   `json:"@context"`
	Type          string   `json:"@type"`
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	Description   string   `json:"description"`
	DatePublished string   `json:"datePublished"`
	Offers        []Offers `json:"offers"`
	Review        []struct {
		Type          string `json:"@type"`
		Author        string `json:"author"`
		DatePublished string `json:"datePublished"`
		Name          string `json:"name"`
		ReviewBody    string `json:"reviewBody"`
		ReviewRating  struct {
			Type        string `json:"@type"`
			RatingValue int    `json:"ratingValue"`
		} `json:"reviewRating"`
	} `json:"review"`
	WorkExample []WorkExample `json:"workExample"`
}

func main() {
	blogTitles, err := GetLatestBlogTitles("https://podcasts.apple.com/podcast/hysteria-radio/id440260757")
	// blogTitles, err := GetLatestBlogTitles("https://podcasts.apple.com/us/podcast/oliver-heldens-presents-heldeep-radio/id1236253646")
	if err != nil {
		log.Println(err)
	}

	for _, song := range blogTitles {
		fmt.Println(song)
	}

}

// GetLatestBlogTitles gets the latest blog title headings from the url
// given and returns them as a list.
func GetLatestBlogTitles(url string) (string, error) {

	// Get the HTML
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(doc.Find("title").Text())

	// var myJSON []string
	doc.Find(".tracks__track--podcast .tracks__track__copy").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			t := regexp.MustCompile(`\d+. `) // backticks are used here to contain the expression
			delimiter1 := regexp.MustCompile(` - `)
			v := t.Split(strings.TrimSpace(s.Text()), -1)

			for _, songTitle := range v {
				//fmt.Printf("%q\n", songTitle)
				slice := delimiter1.Split(songTitle, -1)
				if slice[0] != "" {
					fmt.Printf("%q\n", slice[0])
					fmt.Printf("%q\n\n", slice[len(slice)-1])
				}
			}
		}
	})

	//var podcast Podcast
	//json.Unmarshal([]byte(myJSON[0]), &podcast)
	//workExample := podcast.WorkExample
	// a := regexp.MustCompile(`\d\d?\.\s+|\n\s+`) // a single `a`
	// w := workExample[0]
	// fmt.Printf("%q \n", a.Split(w.Description, -1))

	// fmt.Printf("%q\n", a.Split(workExample.Description, -1))
	//fmt.Printf("Name: %s, \nDescription: %s \n\n", podcast.Name, podcast.WorkExample)

	// for _, v := range myJSON {
	// 	json.Unmarshal([]byte(v), &podcast)
	// 	// fmt.Println(v)
	// 	fmt.Printf("Name: %s, \nDescription: %s \n\n", podcast.Name, podcast.WorkExample)
	// }
	// Save each .post-title as a list
	titles := ""
	doc.Find(".ember171199630").Each(func(i int, s *goquery.Selection) {
		titles += "- " + s.Text() + "\n"
	})
	return titles, nil
}
