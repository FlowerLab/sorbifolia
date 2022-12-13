package puresqlite

import (
	"database/sql"
	"testing"
)

func TestSQLiteVersion(t *testing.T) {
	var version string

	db, err := sql.Open(DriverName, ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if db.QueryRow("SELECT sqlite_version()").Scan(&version) != nil {
		t.Fatal(err)
	}

	t.Log(version)
}
