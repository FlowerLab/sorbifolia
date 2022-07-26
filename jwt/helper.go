package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ES256 = jwt.SigningMethodES256
	ES384 = jwt.SigningMethodES384
	ES512 = jwt.SigningMethodES512
	EdDSA = jwt.SigningMethodEdDSA
	HS256 = jwt.SigningMethodHS256
	HS384 = jwt.SigningMethodHS384
	HS512 = jwt.SigningMethodHS512
	RS256 = jwt.SigningMethodRS256
	RS384 = jwt.SigningMethodRS384
	RS512 = jwt.SigningMethodRS512
	PS256 = jwt.SigningMethodPS256
	PS384 = jwt.SigningMethodPS384
	PS512 = jwt.SigningMethodPS512
	None  = jwt.SigningMethodNone
)

func NewClaims[T any](data *T) *Claims[T]                 { return &Claims[T]{Data: data} }
func (c *Claims[T]) SetID(id string) *Claims[T]           { c.ID = id; return c }
func (c *Claims[T]) SetIssuer(iss string) *Claims[T]      { c.Issuer = iss; return c }
func (c *Claims[T]) SetSubject(sub string) *Claims[T]     { c.Subject = sub; return c }
func (c *Claims[T]) SetAudience(aud ...string) *Claims[T] { c.Audience = aud; return c }
func (c *Claims[T]) AddAudience(aud ...string) *Claims[T] {
	c.Audience = append(c.Audience, aud...)
	return c
}

func (c *Claims[T]) SetExpiresAt(t time.Time) *Claims[T] {
	c.ExpiresAt = jwt.NewNumericDate(t)
	return c
}
func (c *Claims[T]) SetNotBefore(t time.Time) *Claims[T] {
	c.ExpiresAt = jwt.NewNumericDate(t)
	return c
}
func (c *Claims[T]) SetIssuedAt(t time.Time) *Claims[T] {
	c.ExpiresAt = jwt.NewNumericDate(t)
	return c
}
