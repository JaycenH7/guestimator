package parser

import (
	"log"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/mrap/guestimator/models"
	"github.com/mrap/guestimator/request"
)

const topWikisURL = "https://en.wikipedia.org/wiki/Wikipedia:Top_25_Report"

func FetchTopWikis() []models.WikiPage {
	urls := fetchTopWikiURLs()
	wikis := make([]models.WikiPage, len(urls))

	for i, u := range urls {
		if title, err := url.QueryUnescape(path.Base(u)); err == nil {
			if wp, _ := request.GetWikiPage(title); wp != nil {
				wikis[i] = *wp
			}
		}
	}

	return wikis
}

func fetchTopWikiURLs() []string {
	urls := make([]string, 0)

	doc, err := goquery.NewDocument(topWikisURL)
	if err != nil {
		log.Fatalln("Couldn't get top wiki urls", err)
		return urls
	}

	doc.Find(".wikitable tr").Each(func(i int, s *goquery.Selection) {
		if u, exists := s.Find("td a").Attr("href"); exists {
			urls = append(urls, u)
		}
	})

	return urls
}
