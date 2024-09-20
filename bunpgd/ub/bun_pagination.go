package ub

import (
	"errors"
)

func (x *Select) Pagination(v interface {
	Pagination() (limit *int, offset *int, err error)
}) *Select {
	if v == nil {
		return x
	}

	limit, offset, err := v.Pagination()

	switch {
	case err != nil:
		return x.Err(err)
	case limit != nil && offset != nil:
		x.tx.Limit(*limit).Offset(*offset)
	case limit != nil:
		x.tx.Limit(*limit)
	case offset != nil:
		x.tx.Offset(*offset)
	}

	return x
}

type Pagination struct {
	Page     int
	PageSize int
}

func (x *Pagination) Pagination() (*int, *int, error) {
	if x == nil {
		return nil, nil, nil
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
		return nil, nil, nil
	case pageSize < 1:
		pageSize = 1
	case pageSize > 1000:
		pageSize = 1000
	}

	if page*pageSize > 100*10000 {
		return nil, nil, errors.New("page too big")
	}

	page = (page - 1) * pageSize // offset
	return &pageSize, &page, nil
}
