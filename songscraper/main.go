package main

import (
	"fmt"
	"log"
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
	GetSongTitles("https://podcasts.apple.com/podcast/hysteria-radio/id440260757")
}

// GetLatestBlogTitles gets the latest blog title headings from the url
// given and returns them as a list.
func GetSongTitles(url string) error{

	// Convert HTML into goquery document
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	fmt.Println(doc.Find("title").Text())

	// var myJSON []string
	doc.Find("a.link.tracks__track__link--block").Each(func(i int, s *goquery.Selection) {
			fmt.Printf("%s\n", s.Text())
			epURL, ok := s.Attr("href")
			if !ok {
				return
			}
			//fmt.Printf("%d: %s\n", i, epURL)
			doc, err = goquery.NewDocument(epURL)
			if err != nil {
				log.Panicln(err)
			}
			doc.Find("p").Each(func(i int, s *goquery.Selection) {
				t := regexp.MustCompile(`\d+\. `)
				v := t.Split(strings.TrimSpace(s.Text()), -1)
				for _, songTitle := range v {
					if songTitle != "" && songTitle != "1 hr" && songTitle != "Listen on Apple Podcasts" {
						fmt.Printf("%q\n", strings.TrimSpace(songTitle))
					}
				}
				//fmt.Printf("%d: %q\n", i, s.Text())				
			})

			// t := regexp.MustCompile(`\d+. `) // backticks are used here to contain the expression
			// // delimiter1 := regexp.MustCompile(` - `)
			// v := t.Split(strings.TrimSpace(s.Text()), -1)

			// for _, songTitle := range v {
			// 	if songTitle != "" {
			// 		fmt.Printf("%q\n", strings.TrimSpace(songTitle))
			// 	}
				
			// 	// slice := delimiter1.Split(songTitle, -1)
			// 	// if slice[0] != "" {
			// 	// 	fmt.Printf("%q\n", slice[0])
			// 	// 	fmt.Printf("%q\n\n", slice[len(slice)-1])
			// 	// }
			// }
	})
	return nil
}
