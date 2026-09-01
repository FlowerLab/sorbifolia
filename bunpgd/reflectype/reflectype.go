package reflectype

import (
	"reflect"
	"uuid"

	"github.com/uptrace/bun"
)

type BunQueryBuilder interface {
	BunQueryBuilder(bun.QueryBuilder) bun.QueryBuilder
}

type FromQueryParameters interface {
	FromQueryParameters([]string) error
}

var (
	MapStringString = reflect.TypeFor[map[string]string]()
	QueryBuilder    = reflect.TypeFor[BunQueryBuilder]()
	FromQS          = reflect.TypeFor[FromQueryParameters]()
	UUID            = reflect.TypeFor[uuid.UUID]()
)
