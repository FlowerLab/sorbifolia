package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWT[T any] struct {
	method                jwt.SigningMethod
	privateKey, publicKey any
	claims                Claims[T]
}

func New[T any](method jwt.SigningMethod, privateKey, publicKey any, claims Claims[T]) (*JWT[T], error) {
	j := &JWT[T]{
		method:     method,
		privateKey: privateKey,
		publicKey:  publicKey,
		claims:     claims,
	}

	if _, err := j.Generate(Claims[T]{
		RegisteredClaims: jwt.RegisteredClaims{Issuer: "Test", ID: "-1"},
		Data:             new(T),
	}); err != nil {
		return nil, fmt.Errorf("%s: %s", ErrGenSignError, err)
	}

	return j, nil
}

type Claims[T any] struct {
	jwt.RegisteredClaims
	Data *T `json:"data,omitempty"`
}

func (j *JWT[T]) Generate(claims Claims[T]) (string, error) {
	return jwt.NewWithClaims(j.method, claims).SignedString(j.privateKey)
}

func (j *JWT[T]) MustGenerate(claims Claims[T]) string {
	sign, err := j.Generate(claims)
	if err != nil {
		panic(err)
	}
	return sign
}

func (j *JWT[T]) Parse(tokenString string) (*Claims[T], error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		j.claims,
		func(token *jwt.Token) (any, error) { return j.publicKey, nil },
	)
	if err != nil {
		return nil, err
	}

	return j.ParseToken(token)
}

func (j *JWT[T]) ParseToken(token *jwt.Token) (*Claims[T], error) {
	if token == nil {
		return nil, ErrIsNil
	}
	if !token.Valid {
		return nil, ErrNotValid
	}

	claims, ok := token.Claims.(*Claims[T])
	if !ok {
		return nil, ErrClaimsTypeMismatch
	}
	return claims, nil
}

var (
	ErrGenSignError       = errors.New("generate sign error")
	ErrIsNil              = errors.New("token is nil")
	ErrNotValid           = errors.New("not valid")
	ErrClaimsTypeMismatch = errors.New("claims type mismatch")
)
