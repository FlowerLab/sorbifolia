package gomod

import (
	"os"
	"testing"

	"go.x2ox.com/sorbifolia/random"
)

func TestPackage_FindModule(t *testing.T) {
	pkg, err := Parse("example.data.json")
	if err != nil {
		t.Error(err)
	}
	for _, v := range pkg {
		var ms []string
		if ms, err = v.FindModule(); err != nil {
			t.Error(err)
		}

		for _, m := range ms {
			if err = v.Output(random.Fast().RandString(10)); err != nil {
				t.Error(m, err)
			}
		}

		if err = v.Clean(); err != nil {
			t.Error(err)
		}
	}
}

func TestParse(t *testing.T) {
	if _, err := Parse("example.data.json"); err != nil {
		t.Error(err)
	}
	if _, err := Parse("data.json"); err == nil {
		t.Error("err")
	}

	filename := random.Fast().RandString(10)

	file, _ := os.Create(filename)
	_, _ = file.WriteString(`{"a":}`)
	_ = file.Close()
	defer func() {
		_ = os.Remove(filename)
	}()
	if _, err := Parse(filename); err == nil {
		t.Error("err")
	}
}

func TestParseModFile(t *testing.T) {
	if _, err := parseModFile(""); err == nil {
		t.Error("err")
	}
	if _, err := parseModFile("example.data.json"); err == nil {
		t.Error("err")
	}
}
