package builder

import (
	"database/sql"
	"net"
	"testing"

	"github.com/uptrace/bun"
	"go.x2ox.com/sorbifolia/bunpgd"
	"go.x2ox.com/sorbifolia/bunpgd/builder/example"
	"google.golang.org/protobuf/types/known/structpb"
)

var db = bun.NewDB(&sql.DB{}, bunpgd.New())

type User struct {
	bun.BaseModel `bun:"table:user"`
	ID            uint64 `json:"id" bun:",pk"`
}

var testQueryData = []struct {
	sql string
	val []func(*testParseFieldA)
}{
	{sql: `SELECT "user"."id" FROM "user" LIMIT 10`},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("id" = '12') LIMIT 10`,
		val: []func(*testParseFieldA){func(a *testParseFieldA) { a.ID = "12" }},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("username" = 'n') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "n"
				a.Username = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" LIMIT 100 OFFSET 300`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				a.Itr = &example.Pagination{Page: 4, PageSize: 100}
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user"`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				a.Itr = &example.Pagination{Page: 4, PageSize: -1}
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" LIMIT 30`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				a.Itr = &example.Pagination{Page: -1, PageSize: 30}
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE (("username" LIKE '%a%') OR ("username" LIKE '%b%'))`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := []string{"a", "b"}
				a.SearchUsername = &v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE (("username" LIKE '%a') OR ("username" LIKE '%b'))`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := []string{"a", "b"}
				a.QueryUser = &v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("data" @> '{"a":{"b":{"c":"aaa"}}}')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v, _ := structpb.NewValue(map[string]any{
					"a": map[string]any{
						"b": map[string]any{
							"c": "aaa",
						},
					},
				})
				a.DataContain = v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("data" ? 'key')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := "key"
				a.DataKey = &v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("data" ?| '["key1","key2","key3"]')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := []string{"key1", "key2", "key3"}
				a.DataContainKey = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("data" ?& '["key1","key2","key3"]')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := []string{"key1", "key2", "key3"}
				a.DataContainAllKey = &v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("addr" && '1.1.1.1')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				v := net.ParseIP("1.1.1.1")
				a.AddrFb = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("addr" && '1.0.0.0/24')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) {
				_, v, _ := net.ParseCIDR("1.0.0.0/24")
				a.AddrFc = v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("key" ~* 'T.*end') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "T.*end"
				a.KeyMatchCase = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("key" !~ 'T.*end') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "T.*end"
				a.KeyNotMatch = &v
			},
		},
	},

	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("key" ~* 'T.*end') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "T.*end"
				a.KeyMatchCase = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("key" !~* 'T.*end') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "T.*end"
				a.KeyNotMatchCase = &v
			},
		},
	},
	{
		sql: `SELECT "user"."id" FROM "user" WHERE ("key" ^@ 'key_') LIMIT 10`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) {
				v := "key_"
				a.KeyPrefix = &v
			},
		},
	},
}

func TestQuery(t *testing.T) {
	for _, data := range testQueryData {
		var val testParseFieldA
		for _, fn := range data.val {
			fn(&val)
		}

		queryBytes, err := Select(db.NewSelect().Model(&User{}), val).AppendQuery(db.Formatter(), nil)
		if err != nil {
			t.Fatal(err)
		}

		if string(queryBytes) != data.sql {
			t.Errorf("\n got: %v\nwant: %v", string(queryBytes), data.sql)
		}
	}
}
