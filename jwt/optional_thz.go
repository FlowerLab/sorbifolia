//go:build thz

package jwt

import (
	"bytes"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	thz "go.x2ox.com/THz"
)

type THz[T any] struct {
	j         *JWT[T]
	abort     bool
	abortFunc func(ctx *thz.Context)
	store     bool
	storeKey  string
}

func (j *JWT[T]) THz(abort bool, abortFunc func(ctx *thz.Context), store bool, storeKey string) *THz[T] {
	g := &THz[T]{j: j, abort: abort, abortFunc: abortFunc, store: store, storeKey: storeKey}
	if g.abort && g.abortFunc == nil {
		g.abortFunc = func(c *thz.Context) { c.Status(401).Abort() }
	}
	if g.store && g.storeKey == "" {
		g.storeKey = "JWT"
	}
	return g
}

func (g *THz[T]) THzMiddleware() thz.Handler {
	return func(c *thz.Context) {
		claims, err := g.Parse(c)
		if err != nil {
			if g.abort {
				g.abortFunc(c)
			}
			return
		}
		if g.store {
			c.Set(g.storeKey, claims.Data)
		}
		c.Next()
	}
}

func (g *THz[T]) Parse(c *thz.Context) (*Claims[T], error) {
	authHeader := c.Request().Header.Peek("Authorization")
	if len(authHeader) < 7 || !bytes.EqualFold(authHeader[:7], []byte("Bearer ")) {
		return nil, errors.New("invalid authorization")
	}

	token, err := new(jwt.Parser).ParseWithClaims(
		string(authHeader[7:]),
		new(Claims[T]),
		func(token *jwt.Token) (interface{}, error) { return g.j.publicKey, nil },
	)
	if err != nil {
		return nil, err
	}

	return g.j.ParseToken(token)
}
