package signater

import (
	"iter"
	"net/url"
	"strconv"
)

// OrderDirection is the sort direction accepted by list endpoints.
type OrderDirection string

// Sort directions accepted by the OrderByDirection query parameter.
const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

// ListParams are the pagination parameters shared by all list endpoints.
// Zero values fall back to the API defaults: page size 10 (max 100), page
// number 1, descending order.
type ListParams struct {
	// PageSize is the number of items per page (1-100).
	PageSize int
	// PageNumber is the 1-based page to fetch.
	PageNumber int
	// OrderByDirection sorts results ascending or descending.
	OrderByDirection OrderDirection
}

// values encodes the parameters as their query-string representation,
// applying API defaults for zero values.
func (p ListParams) values() url.Values {
	size := p.PageSize
	if size <= 0 {
		size = 10
	}
	page := p.PageNumber
	if page <= 0 {
		page = 1
	}
	dir := p.OrderByDirection
	if dir == "" {
		dir = OrderDesc
	}
	v := url.Values{}
	v.Set("PageSize", strconv.Itoa(size))
	v.Set("PageNumber", strconv.Itoa(page))
	v.Set("OrderByDirection", string(dir))
	return v
}

// Pagination describes the position of a Page within the full result set.
type Pagination struct {
	// TotalItems is the total number of items across all pages.
	TotalItems int `json:"totalItems"`
	// PageSize is the requested number of items per page.
	PageSize int `json:"pageSize"`
	// PageNumber is the 1-based number of this page.
	PageNumber int `json:"pageNumber"`
	// PageItems is the number of items actually present in this page.
	PageItems int `json:"pageItems"`
}

// Page is a single page of results from a list endpoint.
type Page[T any] struct {
	// Items holds the page's results.
	Items []T `json:"items"`
	// Pagination locates this page within the full result set.
	Pagination Pagination `json:"pagination"`
}

// autoPaging returns an iterator that yields every item across pages, calling
// fetch with increasing page numbers starting at firstPage (minimum 1). It
// stops after a short page, once totalItems have been covered, when the
// server omits pagination metadata (nothing to safely continue on), or when
// fetch returns an error, which is yielded as the final pair.
func autoPaging[T any](firstPage int, fetch func(pageNumber int) (*Page[T], error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		page := max(firstPage, 1)
		for {
			p, err := fetch(page)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if p == nil {
				return
			}
			for _, item := range p.Items {
				if !yield(item, nil) {
					return
				}
			}
			pg := p.Pagination
			if len(p.Items) == 0 || pg.PageSize <= 0 || pg.PageItems < pg.PageSize {
				return
			}
			if pg.TotalItems > 0 && page*pg.PageSize >= pg.TotalItems {
				return
			}
			page++
		}
	}
}
