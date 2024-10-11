package operator

import (
	"github.com/uptrace/bun/schema"
)

type Operator string

func Parse(v string) Operator {
	return availableOp[v]
}

var _ schema.QueryAppender = Operator("")

func (o Operator) AppendQuery(_ schema.Formatter, b []byte) ([]byte, error) {
	return append(b, o...), nil
}

func (o Operator) String() string { return string(o) }
