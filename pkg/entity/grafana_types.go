package entity

type DashboardSearchResult struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URI   string `json:"uri"`
	URL   string `json:"url"`
}
