//go:build gin

package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
)

type Gin[T any] struct {
	j         *JWT[T]
	abort     bool
	abortFunc func(ctx *gin.Context)
	store     bool
	storeKey  string
}

func (j *JWT[T]) Gin(abort bool, abortFunc func(ctx *gin.Context), store bool, storeKey string) *Gin[T] {
	g := &Gin[T]{j: j, abort: abort, abortFunc: abortFunc, store: store, storeKey: storeKey}
	if g.abort && g.abortFunc == nil {
		g.abortFunc = func(c *gin.Context) { c.AbortWithStatus(http.StatusUnauthorized) }
	}
	if g.store && g.storeKey == "" {
		g.storeKey = "JWT"
	}
	return g
}

func (g *Gin[T]) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

func (g *Gin[T]) Parse(c *gin.Context) (*Claims[T], error) {
	token, err := request.ParseFromRequest(
		c.Request,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) { return g.j.publicKey, nil },
		request.WithClaims(new(Claims[T])),
	)
	if err != nil {
		return nil, err
	}

	return g.j.ParseToken(token)
}
