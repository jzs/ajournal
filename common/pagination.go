package common

import (
	"net/http"
	"strconv"
)

// PaginationArgs Args for pagination
type PaginationArgs struct {
	Limit uint64
	From  string
}

// Pagination represents the pagination constructs
type Pagination struct {
	HasNext bool
	Next    string
	Prev    string
}

// DefaultPagination returns default settings for pagination functionality
func DefaultPagination() PaginationArgs {
	return PaginationArgs{
		Limit: 10,
		From:  "",
	}
}

// ParsePagination parses pagination from request and returns proper pagination args
func ParsePagination(r *http.Request) PaginationArgs {
	slimit := r.FormValue("limit")
	from := r.FormValue("from")

	args := DefaultPagination()
	if slimit != "" {
		if lim, err := strconv.ParseInt(slimit, 10, 64); err == nil {
			args.Limit = uint64(lim)
		}
	}
	args.From = from

	return args
}
