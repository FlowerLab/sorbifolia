package qsparse

import (
	"net"
	"net/netip"
	"net/url"
	"testing"
	"time"

	"go.x2ox.com/sorbifolia/bunpgd/builder/example"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestQS(t *testing.T) {
	var (
		val testParseFieldA
		qs  = "username=" +
			"&name=1&name=2" +
			"&start_at=1729149000000" +
			"&end_at=2000-01-01T01%3A01%3A01.001Z" +
			"&data_contain_key=key1" +
			"&data_contain_key=key2" +
			"&addr_contain=1.0.0.1" +
			"&addr_contain_or_eq=1.1.0.0/24" +
			"&addr_fa=1.0.0.0/24" +
			"&addr_fb=1.0.0.0" +
			"&addr_fc=1.0.0.0/24" +
			"&json=asfhiwuidhuwei" +
			"&data_contain=" + url.QueryEscape(`{"id":1}`) +
			"&itr=1&itr=10"
	)

	uq, _ := url.ParseQuery(qs)
	if err := QS(uq, &val); err != nil {
		t.Fatal(err)
	}
}

type jsonAble struct{}

func (*jsonAble) MarshalJSON() ([]byte, error) { return nil, nil }
func (*jsonAble) UnmarshalJSON([]byte) error   { return nil }

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
