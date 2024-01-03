package sqltype

import (
	"github.com/uptrace/bun/dialect/sqltype"
)

const (
	TimestampTZ = "TIMESTAMPTZ"         // Timestamp with a time zone
	Date        = "DATE"                // Date
	Time        = "TIME"                // Time without a time zone
	TimeTz      = "TIME WITH TIME ZONE" // Time with a time zone
	Interval    = "INTERVAL"            // Time Interval

	INet    = "INET"
	CIDR    = "CIDR"
	MacAddr = "MACADDR"

	SmallSerial = "SMALLSERIAL"
	Serial      = "SERIAL"
	BigSerial   = "BIGSERIAL"

	Char = "CHAR"
	Text = "TEXT"

	Bytea = "BYTEA"

	Boolean = sqltype.Boolean
	HSTORE  = sqltype.HSTORE
	JSON    = sqltype.JSON
	JSONB   = sqltype.JSONB

	SmallInt        = sqltype.SmallInt
	Integer         = sqltype.Integer
	BigInt          = sqltype.BigInt
	Real            = sqltype.Real
	DoublePrecision = sqltype.DoublePrecision
	UnknownType     = "@UnknownType"
)
