package wikipage

import (
	"encoding/json"
	"net/http"
)

type QueryResponse struct {
	Query WikiQuery `json:"query"`
}

func GetWikiPage(title string) (*WikiPage, error) {
	res, err := http.Get(WikiPageUrl(title))
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

	var page WikiPage
	for _, v := range queryRes.Query.Pages {
		page = v
		break
	}

	return &page, nil
}
