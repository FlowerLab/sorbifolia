package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type Info struct {
	ID int `json:"id"`
}

func TestJWT(t *testing.T) {
	t.Parallel()

	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j := New(jwt.SigningMethodEdDSA, rpk, rpk.Public(), Claims[Info]{})

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

	if ts = j.MustGenerate(Claims[Info]{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Abc",
			Subject: "",
			ID:      "-5",
		},
		Data: &Info{ID: 1},
	}); ts == "" {
		t.Error("fail")
	}
}

func TestJWT_MustGenerate(t *testing.T) {
	t.Parallel()

	type TestData struct {
		A, B chan struct{}
	}

	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j := New(jwt.SigningMethodEdDSA, rpk, rpk.Public(), Claims[TestData]{})

	defer func() { _ = recover() }()
	j.MustGenerate(Claims[TestData]{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Abc",
			Subject: "",
			ID:      "-5",
		},
		Data: &TestData{nil, nil},
	})

	t.Error("err")
}

func TestJWT_Parse(t *testing.T) {
	t.Parallel()

	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j := New(jwt.SigningMethodEdDSA, rpk, rpk.Public(), Claims[any]{})

	if _, err := j.Parse("1.2.3"); err == nil {
		t.Error(err)
	}
}

func TestJWT_ParseToken(t *testing.T) {
	t.Parallel()

	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j := New(jwt.SigningMethodEdDSA, rpk, rpk.Public(), Claims[any]{})

	if _, err := j.ParseToken(nil); err == nil {
		t.Error("failed to parse")
	}
	if _, err := j.ParseToken(&jwt.Token{}); err == nil {
		t.Error("failed to parse")
	}
	if _, err := j.ParseToken(&jwt.Token{Valid: true, Claims: nil}); err == nil {
		t.Error("failed to parse")
	}
}
