package sqltype

import (
	"database/sql"
	"encoding/json"
	"net"
	"net/netip"
	"reflect"
	"testing"
	"time"
)

var testDiscoverTypeData = []struct {
	typ reflect.Type
	val string
}{
	{typ: reflect.TypeOf(false), val: Boolean},
	{typ: reflect.TypeOf(uint8(1)), val: SmallInt},
	{typ: reflect.TypeOf(int32(1)), val: Integer},
	{typ: reflect.TypeOf(float32(1)), val: Real},
	{typ: reflect.TypeOf(float64(1)), val: DoublePrecision},
	{typ: reflect.TypeOf(""), val: Text},
	{typ: reflect.TypeOf(map[string]any{}), val: JSONB},
	{typ: reflect.TypeOf([]string{}), val: "TEXT[]"},
	{typ: reflect.TypeOf([]chan string{}), val: UnknownType},
	{typ: reflect.TypeOf(make(chan string)), val: UnknownType},

	{typ: reflect.TypeOf(uint16(1)), val: SmallInt},
	{typ: reflect.TypeOf(int32(1)), val: Integer},
	{typ: reflect.TypeOf(1), val: BigInt},

	{typ: reflect.TypeOf(&encoderJSON{}), val: JSONB},
	{typ: reflect.TypeOf(&encoderText{}), val: Text},
	{typ: reflect.TypeOf(&encoderBinary{}), val: Bytea},
	{typ: reflect.TypeOf(&encoderUnknown{}), val: UnknownType},
	{typ: reflect.TypeOf(encoderUnknown{}), val: UnknownType},

	{typ: reflect.TypeOf(time.Time{}), val: TimestampTZ},
	{typ: reflect.TypeOf(sql.NullBool{}), val: Boolean},
	{typ: reflect.TypeOf(sql.NullFloat64{}), val: DoublePrecision},
	{typ: reflect.TypeOf(sql.NullInt64{}), val: BigInt},
	{typ: reflect.TypeOf(sql.NullInt32{}), val: Integer},
	{typ: reflect.TypeOf(sql.NullInt16{}), val: SmallInt},
	{typ: reflect.TypeOf(sql.NullString{}), val: Text},
	{typ: reflect.TypeOf(json.RawMessage{}), val: JSONB},
	{typ: reflect.TypeOf(netip.Addr{}), val: INet},
	{typ: reflect.TypeOf(netip.Prefix{}), val: CIDR},
	{typ: reflect.TypeOf(net.HardwareAddr{}), val: MacAddr},
}

func TestDiscoverType(t *testing.T) {
	for _, v := range testDiscoverTypeData {
		if val := DiscoverType(v.typ); val != v.val {
			t.Errorf("expected value is %s, but got %s", val, v.val)
		}
	}
}

type (
	encoderJSON    struct{}
	encoderText    struct{}
	encoderBinary  struct{}
	encoderUnknown struct{}
)

func (*encoderJSON) UnmarshalJSON(_ []byte) error     { panic("") }
func (*encoderJSON) MarshalJSON() ([]byte, error)     { panic("") }
func (*encoderText) UnmarshalText(_ []byte) error     { panic("") }
func (*encoderText) MarshalText() ([]byte, error)     { panic("") }
func (*encoderBinary) UnmarshalBinary(_ []byte) error { panic("") }
func (*encoderBinary) MarshalBinary() ([]byte, error) { panic("") }
