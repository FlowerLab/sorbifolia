package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

type Info struct {
	ID int `json:"id"`
}

func TestNew(t *testing.T) {
	gen := Generator{}
	rsaKey, _ := gen.RSA()
	ed25519Key, _ := gen.Ed25519()
	byKey, _ := gen.HMAC()

	if _, err := New(jwt.SigningMethodRS256, rsaKey, rsaKey.PublicKey, Claims[Info]{}); err != nil {
		t.Fatal(err)
	}
	if _, err := New(jwt.SigningMethodRS384, rsaKey, rsaKey.PublicKey, Claims[Info]{}); err != nil {
		t.Fatal(err)
	}
	if _, err := New(jwt.SigningMethodRS512, rsaKey, rsaKey.PublicKey, Claims[Info]{}); err != nil {
		t.Fatal(err)
	}

	if _, err := New(jwt.SigningMethodEdDSA, ed25519Key, ed25519Key.Public(), Claims[Info]{}); err != nil {
		t.Fatal(err)
	}

	if _, err := New(jwt.SigningMethodHS256, byKey, byKey, Claims[Info]{}); err != nil {
		t.Fatal(err)
	}
	if _, err := New(jwt.SigningMethodEdDSA, byKey, byKey, Claims[Info]{}); err == nil {
		t.Fatal(err)
	}

}

func TestJWT_Generate(t *testing.T) {
	gen := Generator{}
	rpk, _ := gen.Ed25519()
	j, _ := New(jwt.SigningMethodEdDSA, rpk, rpk.Public(), Claims[Info]{})

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
	arg := &Claims[Info]{}
	if arg, err = j.Parse(ts); err != nil {
		t.Fatal(err)
	}
	if arg.Data.ID != 1 {
		t.Fatal("err")
	}
}
