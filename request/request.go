package request

import (
	"encoding/json"
	"net/http"

	"github.com/mrap/guestimator/models"
)

type pagesMap map[string]models.Wikipage

func (m pagesMap) isEmpty() bool {
	_, hasKey := m["-1"]
	return hasKey
}

type wikiQuery struct {
	Pages pagesMap `json:"pages"`
}

type queryResponse struct {
	Query wikiQuery `json:"query"`
}

func GetWikipage(title string) (*models.Wikipage, error) {
	return GetWikipageByUrl(models.WikipageUrl(title))
}

func GetWikipageByUrl(url string) (*models.Wikipage, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	var queryRes queryResponse
	if err = decoder.Decode(&queryRes); err != nil {
		return nil, err
	}

	pages := queryRes.Query.Pages
	if pages.isEmpty() {
		return nil, nil
	}

	var page models.Wikipage
	for _, v := range queryRes.Query.Pages {
		page = v
		break
	}

	return &page, nil
}
