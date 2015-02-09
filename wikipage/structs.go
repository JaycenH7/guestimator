package wikipage

type PagesMap map[string]WikiPage

func (m PagesMap) IsEmpty() bool {
	_, hasKey := m["-1"]
	return hasKey
}

type WikiQuery struct {
	Pages PagesMap `json:"pages"`
}

type WikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}
