package entity

type MetricUsageCount struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type AlertingRule struct {
	Name        string `json:"name"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}
