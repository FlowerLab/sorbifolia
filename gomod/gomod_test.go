package gomod

import (
	"os"
	"testing"
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
			if err = v.Output(m); err != nil {
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
	file, _ := os.Create("data.json")
	_, _ = file.WriteString(`{"a":}`)
	_ = file.Close()
	defer func() {
		_ = os.Remove("data.json")
	}()
	if _, err := Parse("data.json"); err == nil {
		t.Error("err")
	}
}

// func main() {
// 	flag.Parse()
//
// 	pkg, err := gomod.Parse(filename)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	for _, v := range pkg {
// 		var ms []string
// 		if ms, err = v.FindModule(); err != nil {
// 			log.Println(err)
// 		}
//
// 		for _, m := range ms {
// 			if err = v.Output(m); err != nil {
// 				log.Println(m, err)
// 			}
// 		}
//
// 		if err = v.Clean(); err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func TestParseModFile(t *testing.T) {
	if _, err := parseModFile(""); err == nil {
		t.Error("err")
	}
	if _, err := parseModFile("example.data.json"); err == nil {
		t.Error("err")
	}
}
