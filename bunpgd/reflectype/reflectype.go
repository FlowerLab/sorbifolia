package reflectype

import (
	"reflect"

	"github.com/uptrace/bun"
)

type BunQueryBuilder interface {
	BunQueryBuilder(bun.QueryBuilder) bun.QueryBuilder
}

var (
	MapStringString = reflect.TypeFor[map[string]string]()
	QueryBuilder    = reflect.TypeFor[BunQueryBuilder]()
)
