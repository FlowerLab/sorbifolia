package query

import (
	"net"
	"net/netip"
	"reflect"
	"strconv"
	"testing"
	"time"

	"go.x2ox.com/sorbifolia/bunpgd/builder/example"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/flag"
	"google.golang.org/protobuf/types/known/structpb"
)

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

var testParseFieldData = []Field{
	{Name: "id", Flag: flag.String, Op: "=", Key: "id"},
	{Name: "username", Flag: flag.String | flag.Pointer, Op: "=", Key: "username"},
	{Name: "password", Flag: flag.String | flag.Slice, Op: "IN", Key: "password"},
	{Name: "name", Flag: flag.Pointer | flag.String | flag.Slice, Op: "IN", Key: "name"},
	{Name: "json", Flag: flag.Pointer | flag.JSON, Op: "=", Key: "json"},
	{Name: "itr", Flag: flag.Pointer | flag.BunQueryItr, Op: "", Key: "itr"},

	{Name: "start_at", Flag: flag.Pointer | flag.JSON, Op: ">=", Key: "at"},
	{Name: "end_at", Flag: flag.Pointer | flag.JSON, Op: "<=", Key: "at"},

	{Name: "search_username", Flag: flag.Pointer | flag.String | flag.Slice, Op: "LIKE", Key: "username"},
	{Name: "query_user", Flag: flag.Pointer | flag.String | flag.Slice, Op: "LIKE", Key: "username", Attr: "L"},

	{Name: "data_contain", Flag: flag.Pointer | flag.JSON, Op: "@>", Key: "data"},
	{Name: "data_key", Flag: flag.Pointer | flag.String, Op: "?", Key: "data"},
	{Name: "data_contain_key", Flag: flag.Pointer | flag.String | flag.Slice, Op: "?|", Key: "data"},
	{Name: "data_contain_all_key", Flag: flag.Pointer | flag.String | flag.Slice, Op: "?&", Key: "data"},

	{Name: "addr_contain", Flag: flag.Pointer | flag.IP, Op: ">>", Key: "addr"},
	{Name: "addr_contain_or_eq", Flag: flag.Pointer | flag.IP, Op: ">>=", Key: "addr"},
	{Name: "addr_contain_by", Flag: flag.Pointer | flag.IP, Op: "<<", Key: "addr"},
	{Name: "addr_contain_by_or_eq", Flag: flag.Pointer | flag.IP, Op: "<<=", Key: "addr"},
	{Name: "addr_overlap", Flag: flag.Pointer | flag.IP, Op: "&&", Key: "addr"},

	{Name: "addr_fa", Flag: flag.Pointer | flag.IP, Op: "&&", Key: "addr"},
	{Name: "addr_fb", Flag: flag.Pointer | flag.IP, Op: "&&", Key: "addr"},
	{Name: "addr_fc", Flag: flag.Pointer | flag.IP, Op: "&&", Key: "addr"},

	{Name: "key_prefix", Flag: flag.Pointer | flag.String, Op: "^@", Key: "key"},
	{Name: "key_match", Flag: flag.Pointer | flag.String, Op: "~", Key: "key"},
	{Name: "key_not_match", Flag: flag.Pointer | flag.String, Op: "!~", Key: "key"},
	{Name: "key_match_case", Flag: flag.Pointer | flag.String, Op: "~*", Key: "key"},
	{Name: "key_not_match_case", Flag: flag.Pointer | flag.String, Op: "!~*", Key: "key"},
}

func TestParseField(t *testing.T) {
	rt := reflect.TypeFor[testParseFieldA]()

	for i, val := range testParseFieldData {
		sf := rt.Field(i)
		res := ParseField(sf)
		res.Typ = nil

		if !reflect.DeepEqual(res, val) {
			t.Errorf("ParseField(%v) returned %v, wanted %v", sf, res, val)
		}
	}
}

var testParseFieldErrData = []struct {
	panic bool
	null  bool
	v     any
}{
	{panic: false, null: true, v: struct {
		A string `json:"a" query:"-"`
	}{}},
	{panic: false, null: false, v: struct {
		A string `json:"-" query:"-"`
	}{}},
	{panic: false, null: false, v: struct {
		A string `json:"" query:"-"`
	}{}},
	{panic: false, null: false, v: struct {
		A string `json:"" query:"a"`
	}{}},

	{panic: false, null: false, v: struct {
		A string `json:"" query:"a,op:in"`
	}{}},

	{panic: true, null: false, v: struct {
		A string `json:"" query:"a,opx:in"`
	}{}},
	{panic: true, null: false, v: struct {
		A string `json:"" query:"a,op:op"`
	}{}},
	{panic: true, null: false, v: struct {
		A string `json:"" query:"a,opop"`
	}{}},
	{panic: false, null: false, v: struct {
		A string
	}{}},
}

func TestParseFieldErr(t *testing.T) {
	for i, val := range testParseFieldErrData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !val.panic {
						t.Errorf("err")
					}
				}
			}()

			rt := reflect.ValueOf(val.v).Type()

			sf := rt.Field(0)
			res := ParseField(sf)

			if val.null {
				if res.Name != "" {
					t.Errorf("Expected null, got %s", res.Name)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	var v testParseFieldA
	iter := Generate(v)
	iter(func(_ HandleFunc) bool { return false })
}
