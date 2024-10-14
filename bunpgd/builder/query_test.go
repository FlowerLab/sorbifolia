package builder

import (
	"database/sql"
	"net"
	"net/netip"
	"testing"
	"time"

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

type jsonAble struct{}

func (*jsonAble) MarshalJSON() ([]byte, error) { return nil, nil }

type testParseFieldA struct {
	ID       string              `json:"id"`
	Username *string             `json:"username"`
	Password []string            `json:"password"`
	Name     *[]string           `json:"name"`
	JSON     *jsonAble           `json:"json"`
	Itr      *example.Pagination `json:"itr"`

	StartAt *time.Time `json:"start_at" query:"key:at,op:greater_than_or_eq"`
	EndAt   *time.Time `json:"end_at" query:"key:at,op:loe"`

	SearchUsername *[]string `json:"search_username" query:"op:like,key:username"`
	QueryUser      *[]string `json:"query_user" query:"op:like,key:username,attr:L"`

	DataContain       *structpb.Value `json:"data_contain" query:"key:data,op:contain"`
	DataKey           *string         `json:"data_key" query:"key:data,op:exist"`
	DataContainKey    *[]string       `json:"data_contain_key" query:"key:data,op:contain_key"`
	DataContainAllKey *[]string       `json:"data_contain_all_key" query:"key:data,op:contain_all_key"`

	AddrContain       *netip.Addr   `json:"addr_contain" query:"key:addr,op:subnet_contain"`
	AddrContainOrEq   *netip.Prefix `json:"addr_contain_or_eq" query:"key:addr,op:subnet_contain_or_eq"`
	AddrContainBy     *netip.Addr   `json:"addr_contain_by" query:"key:addr,op:subnet_contain_by"`
	AddrContainByOrEq *netip.Prefix `json:"addr_contain_by_or_eq" query:"key:addr,op:subnet_contain_by_or_eq"`
	AddrOverlap       *netip.Addr   `json:"addr_overlap" query:"key:addr,op:subnet_overlap"`

	AddrFa *netip.Prefix `json:"addr_fa" query:"key:addr,op:subnet_overlap"`
	AddrFb *net.IP       `json:"addr_fb" query:"key:addr,op:subnet_overlap"`
	AddrFc *net.IPNet    `json:"addr_fc" query:"key:addr,op:subnet_overlap"`

	KeyPrefix       *string `json:"key_prefix" query:"key:key,op:starts_with"`
	KeyMatch        *string `json:"key_match" query:"key:key,op:regex"`
	KeyNotMatch     *string `json:"key_not_match" query:"key:key,op:not_regex"`
	KeyMatchCase    *string `json:"key_match_case" query:"key:key,op:regex_i"`
	KeyNotMatchCase *string `json:"key_not_match_case" query:"key:key,op:not_regex_i"`
}

var testSelectData = []struct {
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

func TestSelect(t *testing.T) {
	for _, data := range testSelectData {
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

var testUpdateData = []struct {
	sql string
	val []func(*testParseFieldA)
}{
	{
		sql: `UPDATE "user" AS "user" SET a = 'b' WHERE ("id" = '12')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) { a.ID = "12" },
		},
	},
	{
		sql: `UPDATE "user" AS "user" SET a = 'b' WHERE ("id" = '12')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: 100} },
			func(a *testParseFieldA) { a.ID = "12" },
		},
	},
}

func TestUpdate(t *testing.T) {
	for _, data := range testUpdateData {
		var val testParseFieldA
		for _, fn := range data.val {
			fn(&val)
		}

		queryBytes, err := Update(db.NewUpdate().Model(&User{}), val).Set("a = ?", "b").
			AppendQuery(db.Formatter(), nil)
		if err != nil {
			t.Fatal(err)
		}

		if string(queryBytes) != data.sql {
			t.Errorf("\n got: %v\nwant: %v", string(queryBytes), data.sql)
		}
	}
}

var testDeleteData = []struct {
	sql string
	val []func(*testParseFieldA)
}{
	{
		sql: `DELETE FROM "user" AS "user" WHERE ("id" = '12')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: -1} },
			func(a *testParseFieldA) { a.ID = "12" },
		},
	},
	{
		sql: `DELETE FROM "user" AS "user" WHERE ("id" = '12')`,
		val: []func(*testParseFieldA){
			func(a *testParseFieldA) { a.Itr = &example.Pagination{PageSize: 100} },
			func(a *testParseFieldA) { a.ID = "12" },
		},
	},
}

func TestDelete(t *testing.T) {
	for _, data := range testDeleteData {
		var val testParseFieldA
		for _, fn := range data.val {
			fn(&val)
		}

		queryBytes, err := Delete(db.NewDelete().Model(&User{}), val).
			AppendQuery(db.Formatter(), nil)
		if err != nil {
			t.Fatal(err)
		}

		if string(queryBytes) != data.sql {
			t.Errorf("\n got: %v\nwant: %v", string(queryBytes), data.sql)
		}
	}
}
