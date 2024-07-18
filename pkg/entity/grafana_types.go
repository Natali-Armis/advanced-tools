package entity

type DashboardSearchResult struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URI   string `json:"uri"`
	URL   string `json:"url"`
}

type DashboardResponse struct {
	Dashboard Dashboard `json:"dashboard"`
}

type Dashboard struct {
	Title  string  `json:"title"`
	Panels []Panel `json:"panels"`
}

type Panel struct {
	Targets []Target `json:"targets"`
}

type Target struct {
	Expr string `json:"expr"`
}
