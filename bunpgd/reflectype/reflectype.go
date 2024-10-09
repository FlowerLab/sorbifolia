package reflectype

import (
	"reflect"

	"github.com/uptrace/bun"
)

type BunQuery interface {
	BunSelectQuery(*bun.SelectQuery) *bun.SelectQuery
}

var (
	MapStringString = reflect.TypeFor[map[string]string]()

	BunSelectQuery = reflect.TypeFor[BunQuery]()
)
