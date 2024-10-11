package builder

import (
	"reflect"
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

	DataContain    *structpb.Value `json:"data_contain" query:"key:data,op:contain"`
	DataKey        *string         `json:"data_key" query:"key:data,op:exist"`
	DataContainKey *[]string       `json:"data_contain_key" query:"key:data,op:contain_key"`
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
