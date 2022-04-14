package gomod

import (
	"testing"
)

func TestPackage_FindModule(t *testing.T) {
	pkg, err := Parse("example.config.json")
	if err != nil {
		t.Error(err)
	}
	for _, v := range pkg {
		var ms []string
		if ms, err = v.FindModule(); err != nil {
			t.Error(err)
		}

		for _, m := range ms {
			if err = v.Output(m); err != nil {
				t.Error(m, err)
			}
		}

		if err = v.Clean(); err != nil {
			t.Error(err)
		}
	}
}
