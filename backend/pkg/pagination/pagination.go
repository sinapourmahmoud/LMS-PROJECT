package pagination

import (
	"errors"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Params struct {
	Page  int
	Limit int
}

func (p Params) Offset() int { return (p.Page - 1) * p.Limit }

// Page is the response of every paginated endpoint.
type Page[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func NewPage[T any](items []T, params Params, total int) *Page[T] {
	return &Page[T]{Items: items, Page: params.Page, Limit: params.Limit, Total: total}
}

// Parse turns the ?page= and ?limit= query values into Params.
// An empty value means "use the default".
func Parse(pageStr string, limitStr string) (Params, error) {
	params := Params{Page: DefaultPage, Limit: DefaultLimit}

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return Params{}, errors.New("page must be a number greater than 0")
		}
		params.Page = page
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > MaxLimit {
			return Params{}, errors.New("limit must be a number between 1 and " + strconv.Itoa(MaxLimit))
		}
		params.Limit = limit
	}

	return params, nil
}