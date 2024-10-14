package builder

import (
	"github.com/uptrace/bun"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/query"
)

func Select(q *bun.SelectQuery, v any) *bun.SelectQuery {
	for val := range query.Generate(v) {
		q.ApplyQueryBuilder(val)
	}
	return q
}

func Update(q *bun.UpdateQuery, v any) *bun.UpdateQuery {
	for val := range query.Generate(v) {
		q.ApplyQueryBuilder(val)
	}
	return q
}

func Delete(q *bun.DeleteQuery, v any) *bun.DeleteQuery {
	for val := range query.Generate(v) {
		q.ApplyQueryBuilder(val)
	}
	return q
}
