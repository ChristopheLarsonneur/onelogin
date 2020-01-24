package api

// ResultPagination is the Pagination section for some queries
// as described in https://developers.onelogin.com/api-docs/1/getting-started/using-query-parameters#pagination
type ResultPagination struct {
	BeforeCursor string `json:"before_cursor"`
	AfterCursor string `json:"after_cursor"`
}
