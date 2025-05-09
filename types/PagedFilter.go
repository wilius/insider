package types

import (
	"net/url"
	"strconv"
)

type PagedFilter struct {
	Page int
	Size int
}

func (r *PagedFilter) CalculateOffset() int {
	return r.Size * (r.Page - 1)
}

func (r *PagedFilter) CalculateLimit() int {
	return r.Size + 1
}

func ParseQueryForPageFilter(values *url.Values) *PagedFilter {
	page, _ := strconv.Atoi(values.Get("page"))
	size, _ := strconv.Atoi(values.Get("size"))

	return NewPagedFilter(page, size)

}

func NewPagedFilter(page int, size int) *PagedFilter {
	if page <= 0 {
		page = 1
	}

	if size == 0 {
		size = 20
	}

	return &PagedFilter{
		Page: page,
		Size: size,
	}
}
