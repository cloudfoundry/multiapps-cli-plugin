package models

type PaginatedResponse[T any] struct {
	Pagination Pagination `json:"pagination"`
	Resources  []T        `json:"resources"`
}

type Pagination struct {
	TotalPages int64  `json:"total_pages"`
	NextPage   string `jsonry:"next.href"`
}

type PaginatedResponseWithIncluded[T any, Included any] struct {
	Pagination Pagination `json:"pagination"`
	Resources  []T        `json:"resources"`
	Included   Included   `json:"included"`
}
