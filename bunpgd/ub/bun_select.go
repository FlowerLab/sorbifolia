package ub

import (
	"errors"
)

type PaginationItr interface {
	Pagination() (limit int, offset int, err error)
}

type SortItr interface {
	Sort() []string
}

func (x *Select) Pagination(v PaginationItr) *Select {
	if v == nil {
		return x
	}

	limit, offset, err := v.Pagination()

	switch {
	case err != nil:
		return x.Err(err)
	case limit > 0 && offset > 0:
		x.tx.Limit(limit).Offset(offset)
	case limit > 0:
		x.tx.Limit(limit)
	case offset > 0:
		x.tx.Offset(offset)
	}

	return x
}

type Pagination struct {
	Page     int
	PageSize int
}

func (x *Pagination) Pagination() (int, int, error) {
	if x == nil {
		return -1, -1, nil
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
		return -1, -1, nil
	case pageSize < 1:
		pageSize = 1
	case pageSize > 1000:
		pageSize = 1000
	}

	if page*pageSize > 100*10000 {
		return -1, -1, errors.New("page too big")
	}
	return pageSize, (page - 1) * pageSize, nil
}

func (x *Select) Sort(v SortItr) *Select {
	if v == nil {
		return x
	}

	if arr := v.Sort(); len(arr) > 0 {
		x.tx.Order(arr...)
	}
	return x
}

type Sort []SortItem

type SortItem struct {
	Key    string `json:"key"`
	IsDesc bool   `json:"is_desc"`
}

func (x *Sort) Sort() []string {
	if x == nil {
		return nil
	}

	arr := make([]string, 0, len(*x))
	for _, v := range *x {
		if v.Key == "" {
			continue
		}

		if v.IsDesc {
			arr = append(arr, v.Key+" DESC")
		} else {
			arr = append(arr, v.Key)
		}
	}
	if len(arr) == 0 {
		return nil
	}

	return arr
}
