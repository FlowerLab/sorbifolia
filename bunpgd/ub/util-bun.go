//go:generate go run gen_bun.go
//go:generate go run gen_extend.go
package ub

import (
	"sync/atomic"

	"github.com/uptrace/bun"
)

type U struct{ *bun.DB }

var defaultH atomic.Pointer[U]

func Get() *U           { return defaultH.Load() }
func Set(db *bun.DB)    { defaultH.Store(New(db)) }
func New(db *bun.DB) *U { return &U{db} }
