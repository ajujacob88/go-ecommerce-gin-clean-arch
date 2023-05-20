package model

type QueryParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Query  string `json:"query"`
	Filter string `json:"filter"`
	SortBy string `json:"sort_by"`
	//sortDesc bool   `json:"sort_desc"`
	SortDesc bool `json:"sort_desc"`
}
