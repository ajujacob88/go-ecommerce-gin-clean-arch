package common

//model structs used for both input and output

type QueryParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Query    string `json:"query"`
	Filter   string `json:"filter"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}
