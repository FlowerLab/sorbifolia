package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/attribute"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/flag"
	op "go.x2ox.com/sorbifolia/bunpgd/builder/internal/operator"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

func ParseField(sf reflect.StructField) (field Field) {
	field.Typ = sf.Type

	query := sf.Tag.Get("query")
	if query == "-" {
		return
	}

	var arr []string
	if query != "" {
		if arr = strings.Split(query, ","); !strings.Contains(arr[0], ":") {
			field.Name = arr[0]
			arr = arr[1:]
		}
	}

	for _, kv := range arr {
		k, v, has := strings.Cut(kv, ":")
		if !has {
			panic(fmt.Sprintf("parse %s: %s(%s) error, notfound delimiter %s", sf.PkgPath, sf.Name, sf.Type, kv))
		}

		switch k {
		case "op":
			if field.Op = op.Parse(v); field.Op == op.Unknown {
				panic(fmt.Sprintf("parse %s: %s(%s) error, unknown operator %s", sf.PkgPath, sf.Name, sf.Type, k))
			}
		case "attr":
			field.Attr = attribute.Attribute(v)
		case "key":
			field.Key = bun.Ident(v)
		default:
			panic(fmt.Sprintf("parse %s: %s(%s) error, unknown key %s", sf.PkgPath, sf.Name, sf.Type, k))
		}
	}

	if field.Name == "" {
		if val := sf.Tag.Get("json"); val != "" && val != "-" {
			field.Name, _, _ = strings.Cut(val, ",") // ignore opt
		}
		if field.Name == "" || field.Name == "-" {
			field.Name = sf.Name
		}
	}

	if field.Key == "" {
		field.Key = bun.Ident(field.Name)
	}

	if field.Flag = flag.From(field.Typ); field.Flag.And(reflectype.QueryBuilder) {
		field.Op = op.Unknown
		return
	}

	if field.Op == op.Unknown {
		typ := field.Typ

		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}

		switch typ.Kind() {
		case reflect.Slice:
			field.Op = op.In
		default:
			field.Op = op.Equal
		}
	}

	return
}
