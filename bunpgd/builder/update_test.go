package builder

import (
	"errors"
	"fmt"
	"testing"
)

type testOptionalUpdate struct {
	ID      uint64     `json:"id"`
	SubID   uint64     `json:"sub_id"`
	XID     *uint      `json:"-"`
	Enable  bool       `json:"enable"`
	Disable *bool      `json:"disable"`
	JSON    *jsonAbleX `json:"json"`
	User
}

type jsonAbleX struct{}

func (*jsonAbleX) MarshalJSON() ([]byte, error) { return []byte("2"), nil }

func TestOptionalUpdate(t *testing.T) {
	{
		queryBytes, err := OptionalUpdate(db.NewUpdate().Model(&User{}).Where("1=1"), &testOptionalUpdate{
			ID:     0,
			SubID:  0,
			XID:    nil,
			Enable: false,
			JSON:   &jsonAbleX{},
		}, "sub_id").AppendQuery(db.Formatter(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(queryBytes) != `UPDATE "user" AS "user" SET "id" = 0, "enable" = FALSE, "json" = '2' WHERE (1=1)` {
			t.Fatal(fmt.Errorf("unexpected query: %s", string(queryBytes)))
		}
	}

	{
		_, err := OptionalUpdate(db.NewUpdate().Model(&User{}), map[string]any{}).AppendQuery(db.Formatter(), nil)
		if err == nil {
			t.Fatal(errors.New("expected error"))
		}
	}
}

func TestOptionalForceUpdate(t *testing.T) {
	{
		queryBytes, err := OptionalForceUpdate(db.NewUpdate().Model(&User{}).Where("1=1"), &testOptionalUpdate{
			ID:     0,
			SubID:  0,
			XID:    nil,
			Enable: false,
			JSON:   &jsonAbleX{},
		}, []string{"disable"}, []string{"sub_id"}).AppendQuery(db.Formatter(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(queryBytes) != `UPDATE "user" AS "user" SET "id" = 0, "enable" = FALSE, "disable" = NULL, "json" = '2' WHERE (1=1)` {
			t.Fatal(fmt.Errorf("unexpected query: %s", string(queryBytes)))
		}
	}

	{
		_, err := OptionalForceUpdate(db.NewUpdate().Model(&User{}), map[string]any{}, nil, nil).AppendQuery(db.Formatter(), nil)
		if err == nil {
			t.Fatal(errors.New("expected error"))
		}
	}
}
