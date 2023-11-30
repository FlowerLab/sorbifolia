package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewClaims(t *testing.T) {
	t.Parallel()

	type TestData struct {
		Data string `json:"data"`
	}
	nt := time.Now().Truncate(jwt.TimePrecision)
	claims := NewClaims(new(TestData))

	if !claims.SetIssuedAt(nt).IssuedAt.Equal(nt) {
		t.Error("fail")
	}
	if !claims.SetExpiresAt(nt).ExpiresAt.Equal(nt) {
		t.Error("fail")
	}
	if !claims.SetNotBefore(nt).NotBefore.Equal(nt) {
		t.Error("fail")
	}

	if claims.SetID("a").ID != "a" {
		t.Error("fail")
	}
	if claims.SetIssuer("a").Issuer != "a" {
		t.Error("fail")
	}
	if claims.SetSubject("a").Subject != "a" {
		t.Error("fail")
	}
	if claims.SetAudience("a").Audience[0] != "a" {
		t.Error("fail")
	}
	if claims.AddAudience("b").Audience[1] != "b" {
		t.Error("fail")
	}
}
