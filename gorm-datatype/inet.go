package datatype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/netip"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	INetAddr   netip.Addr
	INetPrefix netip.Prefix
)

func (*INetAddr) GormDataType() string   { return "INetAddr" }
func (*INetPrefix) GormDataType() string { return "INetPrefix" }

func (*INetPrefix) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	return isPostgres(db, "inet")
}

func (*INetAddr) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	return isPostgres(db, "inet")
}

// Scan implements the sql.Scanner interface.
func (a *INetAddr) Scan(value any) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("failed to unmarshal ip value: %v", value)
	}

	addr, err := netip.ParseAddr(str)
	if err == nil {
		*a = INetAddr(addr)
	}

	return err
}

// Value implements the driver.Valuer interface.
func (a INetAddr) Value() (driver.Value, error) {
	na := netip.Addr(a)
	if na.IsValid() {
		return na.String(), nil
	}
	return "", errors.New("invalid ip address")
}

// Scan implements the sql.Scanner interface.
func (a *INetPrefix) Scan(value any) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("failed to unmarshal ip value: %v", value)
	}

	addr, err := netip.ParsePrefix(str)
	if err == nil {
		*a = INetPrefix(addr)
	}

	return err
}

// Value implements the driver.Valuer interface.
func (a INetPrefix) Value() (driver.Value, error) {
	np := netip.Prefix(a)
	if np.IsValid() {
		return np.String(), nil
	}
	return "", errors.New("invalid ip prefix")
}
