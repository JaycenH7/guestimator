package request

import (
	"encoding/json"
	"net/http"

	"github.com/mrap/guestimator/models"
)

type PagesMap map[string]models.WikiPage

func (m PagesMap) IsEmpty() bool {
	_, hasKey := m["-1"]
	return hasKey
}

type WikiQuery struct {
	Pages PagesMap `json:"pages"`
}

type QueryResponse struct {
	Query WikiQuery `json:"query"`
}

func GetWikiPage(title string) (*models.WikiPage, error) {
	res, err := http.Get(models.WikiPageUrl(title))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	var queryRes QueryResponse
	if err = decoder.Decode(&queryRes); err != nil {
		return nil, err
	}

	pages := queryRes.Query.Pages
	if pages.IsEmpty() {
		return nil, nil
	}

	var page models.WikiPage
	for _, v := range queryRes.Query.Pages {
		page = v
		break
	}

	return &page, nil
}
