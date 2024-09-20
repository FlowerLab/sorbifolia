//go:build gen_bun

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"reflect"
	"text/template"

	"github.com/uptrace/bun"
)

func main() {
	buf := new(bytes.Buffer)
	writeHead(buf)
	writeType(buf)

	writeFunc(buf, "Select", (&bun.DB{DB: &sql.DB{}}).NewSelect())
	writeFunc(buf, "Insert", (&bun.DB{DB: &sql.DB{}}).NewInsert())
	writeFunc(buf, "Update", (&bun.DB{DB: &sql.DB{}}).NewUpdate())
	writeFunc(buf, "Delete", (&bun.DB{DB: &sql.DB{}}).NewDelete())

	os.WriteFile("bun_gen.go", buf.Bytes(), 0644)
}

func _pt(req, res string) *template.Template {
	const format = "func (x *{{.Query}}) {{.Name}}(%s) *{{.Query}} { x.tx.{{.Name}}(%s); return x }"
	return template.Must(template.New("").Parse(fmt.Sprintf(format, req, res)))
}

var funcMap = map[string]*template.Template{
	"Model": _pt("model any", "model"),

	"Offset": _pt("n int", "n"),
	"Limit":  _pt("n int", "n"),

	"Err": _pt("err error", "err"),

	"UseIndexForJoin":       _pt("args ...string", "args..."),
	"UseIndexForGroupBy":    _pt("args ...string", "args..."),
	"ForceIndex":            _pt("args ...string", "args..."),
	"Table":                 _pt("args ...string", "args..."),
	"Order":                 _pt("args ...string", "args..."),
	"ForceIndexForJoin":     _pt("args ...string", "args..."),
	"Group":                 _pt("args ...string", "args..."),
	"ForceIndexForOrderBy":  _pt("args ...string", "args..."),
	"IgnoreIndexForGroupBy": _pt("args ...string", "args..."),
	"ForceIndexForGroupBy":  _pt("args ...string", "args..."),
	"IgnoreIndexForJoin":    _pt("args ...string", "args..."),
	"WherePK":               _pt("args ...string", "args..."),
	"IgnoreIndexForOrderBy": _pt("args ...string", "args..."),
	"Column":                _pt("args ...string", "args..."),
	"ExcludeColumn":         _pt("args ...string", "args..."),
	"UseIndexForOrderBy":    _pt("args ...string", "args..."),
	"UseIndex":              _pt("args ...string", "args..."),
	"IgnoreIndex":           _pt("args ...string", "args..."),

	"Join":           _pt("s string, args ...any", "s, args..."),
	"GroupExpr":      _pt("s string, args ...any", "s, args..."),
	"DistinctOn":     _pt("s string, args ...any", "s, args..."),
	"TableExpr":      _pt("s string, args ...any", "s, args..."),
	"JoinOnOr":       _pt("s string, args ...any", "s, args..."),
	"ModelTableExpr": _pt("s string, args ...any", "s, args..."),
	"Having":         _pt("s string, args ...any", "s, args..."),
	"ColumnExpr":     _pt("s string, args ...any", "s, args..."),
	"JoinOn":         _pt("s string, args ...any", "s, args..."),
	"Where":          _pt("s string, args ...any", "s, args..."),
	"OrderExpr":      _pt("s string, args ...any", "s, args..."),
	"WhereOr":        _pt("s string, args ...any", "s, args..."),
	"For":            _pt("s string, args ...any", "s, args..."),
	"Returning":      _pt("s string, args ...any", "s, args..."),
	"Set":            _pt("s string, args ...any", "s, args..."),
	"On":             _pt("s string, args ...any", "s, args..."),

	"UnionAll":  _pt("q *bun.{{.Query}}Query", "q"),
	"Except":    _pt("q *bun.{{.Query}}Query", "q"),
	"ExceptAll": _pt("q *bun.{{.Query}}Query", "q"),

	"Apply": _pt("fn func(*bun.{{.Query}}Query) *bun.{{.Query}}Query", "fn"),

	"Relation":   _pt("s string, fn ...func(*bun.{{.Query}}Query) *bun.{{.Query}}Query", "s, fn..."),
	"WhereGroup": _pt("s string, fn func(*bun.{{.Query}}Query) *bun.{{.Query}}Query", "s, fn"),

	"Union":        _pt("q *bun.{{.Query}}Query", "q"),
	"Intersect":    _pt("q *bun.{{.Query}}Query", "q"),
	"IntersectAll": _pt("q *bun.{{.Query}}Query", "q"),

	"WithRecursive": _pt("name string, query schema.QueryAppender", "name, query"),
	"With":          _pt("name string, query schema.QueryAppender", "name, query"),

	"ApplyQueryBuilder": _pt("fn func(bun.QueryBuilder) bun.QueryBuilder", "fn"),

	"Conn": _pt("db bun.IConn", "db"),

	"Distinct":            _pt("", ""),
	"WhereDeleted":        _pt("", ""),
	"WhereAllWithDeleted": _pt("", ""),
	"Replace":             _pt("", ""),
	"Ignore":              _pt("", ""),
	"Bulk":                _pt("", ""),
	"OmitZero":            _pt("", ""),
	"ForceDelete":         _pt("", ""),

	"Value":     _pt("column string, query string, args ...any", "column, query, args..."),
	"SetColumn": _pt("column string, query string, args ...any", "column, query, args..."),
}

type Data struct {
	Query string
	Name  string
}

func writeFunc(w io.Writer, query string, t any) {
	selectMap := getTypeMethod(reflect.TypeOf(t))

	for _, v := range selectMap {
		switch v.Name {
		case "NewSelect", "NewInsert", "NewUpdate", "NewDelete":
			continue
		default:
		}

		format, has := funcMap[v.Name]
		if !has {
			panic(fmt.Sprintf("%s has not impl", v.Name))
		}

		if err := format.Execute(w, Data{Query: query, Name: v.Name}); err != nil {
			panic(err)
		}
		_, _ = w.Write([]byte{'\n'})
	}

	_, _ = w.Write([]byte{'\n'})
}

func writeType(w io.Writer) {
	_, _ = w.Write([]byte(`type (
	Select struct { util *U; tx *bun.SelectQuery }
	Update struct { util *U; tx *bun.UpdateQuery }
	Insert struct { util *U; tx *bun.InsertQuery }
	Delete struct { util *U; tx *bun.DeleteQuery }
)

func (u *U) Select() *Select { return &Select{util: u, tx: u.NewSelect()} }
func (u *U) Update() *Update { return &Update{util: u, tx: u.NewUpdate()} }
func (u *U) Insert() *Insert { return &Insert{util: u, tx: u.NewInsert()} }
func (u *U) Delete() *Delete { return &Delete{util: u, tx: u.NewDelete()} }

func (x *Select) RawQ() *bun.SelectQuery { return x.tx }
func (x *Update) RawQ() *bun.UpdateQuery { return x.tx }
func (x *Insert) RawQ() *bun.InsertQuery { return x.tx }
func (x *Delete) RawQ() *bun.DeleteQuery { return x.tx }

`))
}

func writeHead(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf(`// Code generated by bunpgd/ub. DO NOT EDIT.
// versions: bun %s
// source: bunpgd/ub/gen_bun.go

package ub

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

`, bun.Version())))
}

func getTypeMethod(rt reflect.Type) map[string]reflect.Method {
	m := make(map[string]reflect.Method)
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)

		returnType := method.Type.Out(0)

		if rt == returnType {
			m[method.Name] = method
		}
	}
	return m
}
