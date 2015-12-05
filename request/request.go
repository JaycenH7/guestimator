package request

import (
	"encoding/json"
	"net/http"

	"github.com/mrap/guestimator/models"
)

type pagesMap map[string]models.WikiPage

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

func GetWikiPage(title string) (*models.WikiPage, error) {
	res, err := http.Get(models.WikiPageUrl(title))
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

	var page models.WikiPage
	for _, v := range queryRes.Query.Pages {
		page = v
		break
	}

	return &page, nil
}
