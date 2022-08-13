package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

type Info struct {
	ID int `json:"id"`
}

func TestJWT(t *testing.T) {
	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j := New(EdDSA, rpk, rpk.Public(), Claims[Info]{})

	ts, err := j.Generate(Claims[Info]{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Abc",
			Subject: "",
			ID:      "-5",
		},
		Data: &Info{ID: 1},
	})

	if err != nil {
		t.Fatal(err)
	}
	var arg *Claims[Info]
	if arg, err = j.Parse(ts); err != nil {
		t.Fatal(err)
	}
	if arg.Data.ID != 1 {
		t.Fatal("err")
	}
}
