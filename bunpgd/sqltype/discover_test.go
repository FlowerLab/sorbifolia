package sqltype

import (
	"database/sql"
	"encoding/json"
	"net"
	"net/netip"
	"reflect"
	"testing"
	"time"
	"uuid"
)

var testDiscoverTypeData = []struct {
	typ reflect.Type
	val string
}{
	{typ: reflect.TypeOf(false), val: Boolean},
	{typ: reflect.TypeFor[uint8](), val: SmallInt},
	{typ: reflect.TypeFor[int32](), val: Integer},
	{typ: reflect.TypeFor[float32](), val: Real},
	{typ: reflect.TypeFor[float64](), val: DoublePrecision},
	{typ: reflect.TypeFor[string](), val: Text},
	{typ: reflect.TypeFor[map[string]any](), val: JSONB},
	{typ: reflect.TypeFor[[]string](), val: "TEXT[]"},
	{typ: reflect.TypeFor[[]chan string](), val: UnknownType},
	{typ: reflect.TypeFor[chan string](), val: UnknownType},

	{typ: reflect.TypeFor[uint16](), val: SmallInt},
	{typ: reflect.TypeFor[int32](), val: Integer},
	{typ: reflect.TypeFor[int](), val: BigInt},

	{typ: reflect.TypeFor[*encoderJSON](), val: JSONB},
	{typ: reflect.TypeFor[*encoderText](), val: Text},
	{typ: reflect.TypeFor[*encoderBinary](), val: Bytea},
	{typ: reflect.TypeFor[*encoderUnknown](), val: UnknownType},
	{typ: reflect.TypeFor[encoderUnknown](), val: UnknownType},

	{typ: reflect.TypeFor[time.Time](), val: TimestampTZ},
	{typ: reflect.TypeFor[sql.NullBool](), val: Boolean},
	{typ: reflect.TypeFor[sql.NullFloat64](), val: DoublePrecision},
	{typ: reflect.TypeFor[sql.NullInt64](), val: BigInt},
	{typ: reflect.TypeFor[sql.NullInt32](), val: Integer},
	{typ: reflect.TypeFor[sql.NullInt16](), val: SmallInt},
	{typ: reflect.TypeFor[sql.NullString](), val: Text},
	{typ: reflect.TypeFor[json.RawMessage](), val: JSONB},
	{typ: reflect.TypeFor[netip.Addr](), val: INet},
	{typ: reflect.TypeFor[netip.Prefix](), val: CIDR},
	{typ: reflect.TypeFor[net.HardwareAddr](), val: MacAddr},
	{typ: reflect.TypeFor[uuid.UUID](), val: UUID},
	{typ: reflect.TypeFor[[]uuid.UUID](), val: "UUID[]"},
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
