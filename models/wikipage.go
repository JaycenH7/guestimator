package models

import (
	"log"

	"github.com/mrap/stringutil"
)

const (
	BaseQueryUrl = "http://en.wikipedia.org/w/api.php?action=query&prop=extracts&format=json&exintro=&explaintext=&titles="
)

type WikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

func WikiPageUrl(title string) string {
	encoded, err := stringutil.UrlEncoded(title)
	if err != nil {
		log.Println("Error making wiki page url", err)
	}
	return BaseQueryUrl + encoded
}
