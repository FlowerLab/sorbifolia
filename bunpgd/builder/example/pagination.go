package example

import (
	"errors"
	"strconv"

	"github.com/uptrace/bun"
)

type Pagination struct {
	Page     int
	PageSize int
}

func (x *Pagination) CalcLimitOffset() (limit, offset int) {
	if x == nil {
		return 10, 0 // by default
	}

	var (
		page     = x.Page
		pageSize = x.PageSize
	)

	if page < 1 {
		page = 1
	}

	switch {
	case pageSize == -1:
		return -1, -1
	case pageSize < 1:
		pageSize = 1
	case pageSize > 1000:
		pageSize = 1000
	}

	return pageSize, (page - 1) * pageSize
}

func (x *Pagination) BunQueryBuilder(q bun.QueryBuilder) bun.QueryBuilder {
	limit, offset := x.CalcLimitOffset()

	switch q := q.Unwrap().(type) {
	case *bun.SelectQuery:
		switch {
		case limit+offset > 100*10000:
			q.Err(errors.New("page too big"))
		case limit > 0 && offset > 0:
			q.Limit(limit).Offset(offset)
		case limit > 0:
			q.Limit(limit)
		case offset > 0:
			q.Offset(offset)
		}
	default:
	}

	return q
}

type FromQueryParameters interface {
	FromQueryParameters([]string) error
}

func (x *Pagination) FromQueryParameters(arr []string) (err error) {
	if x == nil {
		return errors.New("x is nil")
	}

	if len(arr) > 0 {
		x.Page, err = strconv.Atoi(arr[0])
	}
	if err != nil && len(arr) > 1 {
		x.PageSize, err = strconv.Atoi(arr[1])
	}
	return
}
