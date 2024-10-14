package operator

import (
	"testing"

	"github.com/uptrace/bun/schema"
)

var testParseData = []struct {
	str string
	op  Operator
}{
	{str: "in", op: In},
	{str: "In", op: Unknown},
}

func TestParse(t *testing.T) {
	for _, tt := range testParseData {
		if v := Parse(tt.str); v != tt.op {
			t.Errorf("Parse(%q) = %v, want %v", tt.str, v, tt.op)
		}
	}
}

func TestOperator_AppendQuery(t *testing.T) {
	b, err := In.AppendQuery(schema.NewNopFormatter(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != len(In) {
		t.Errorf("len(In) = %d, want %d", len(b), len(In))
	}

	buf := make([]byte, 0, 32)

	b, err = In.AppendQuery(schema.NewNopFormatter(), buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != len(In) {
		t.Errorf("len(In) = %d, want %d", len(b), len(In))
	}
	if cap(buf) != 32 {
		t.Errorf("cap(buf) = %d, want %d", cap(buf), 32)
	}
	if string(b) != In.String() {
		t.Errorf("String() = %q, want %q", string(b), In.String())
	}
}
