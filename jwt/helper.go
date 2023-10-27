package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	c.NotBefore = jwt.NewNumericDate(t)
	return c
}
func (c *Claims[T]) SetIssuedAt(t time.Time) *Claims[T] {
	c.IssuedAt = jwt.NewNumericDate(t)
	return c
}
