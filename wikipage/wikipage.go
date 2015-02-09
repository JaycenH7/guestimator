package wikipage

import (
	"log"

	"github.com/mrap/stringutil"
)

const (
	BaseQueryUrl = "http://en.wikipedia.org/w/api.php?action=query&prop=extracts&format=json&exintro=&explaintext=&titles="
)

func WikiPageUrl(title string) string {
	encoded, err := stringutil.UrlEncoded(title)
	if err != nil {
		log.Println("Error making wiki page url", err)
	}
	return BaseQueryUrl + encoded
}
