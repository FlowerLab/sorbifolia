package bunpgd

import (
	"database/sql"
	"strings"
	"testing"
	"uuid"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

func TestDialectUUIDColumns(t *testing.T) {
	// Bun 1.3 discovers stdlib UUIDs as VARCHAR; our dialect must keep UUID.
	type model struct {
		ID  uuid.UUID
		IDs []uuid.UUID `bun:"ids"`
	}
	db := bun.NewDB(new(sql.DB), New())
	query := db.NewCreateTable().Model((*model)(nil)).String()
	for _, column := range []string{`"id" UUID`, `"ids" UUID[]`} {
		if !strings.Contains(query, column) {
			t.Errorf("expected %s in %s", column, query)
		}
	}
}

func TestDialectStringFormatting(t *testing.T) {
	gen := schema.NewQueryGen(New())
	for _, tt := range []struct {
		name, input, want string
	}{
		{"quote", "it's", "SELECT 'it''s'"},
		{"nul", "before\x00after", "SELECT ?!(bun: string contains a NUL byte (0x00))"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := gen.FormatQuery("SELECT ?", tt.input); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
