package gomod

import (
	"bytes"
	"testing"

	"go.x2ox.com/sorbifolia/random"
)

func TestPackageData_Write(t *testing.T) {
	var (
		pd = packageData{
			Main:    "a",
			PkgName: random.SafeRand{}.RandString(10),
			Repo:    "c",
			Branch:  "d",
			ReadMe:  "e",
		}
		buf = &bytes.Buffer{}
	)

	if err := pd.Write(buf); err != nil {
		t.Error(err)
	}
	if buf.Len() == 0 {
		t.Error("no data")
	}

	if err := pd.WriteToFile(); err != nil {
		t.Error(err)
	}

	if err := pd.WriteToFile(); err == nil {
		t.Error("err") // file already exists
	}
}
