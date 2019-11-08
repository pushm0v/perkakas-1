package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/structs"
	"github.com/kitabisa/perkakas/v2/token/jwt"
)

func NewJWT(hctx phttp.HttpHandlerContext, signKey []byte) func(next http.Handler) http.Handler {
	jwtt := jwt.NewJWT(signKey)
	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			match, err := regexp.MatchString("^Bearer .+", authorization)
			if err != nil || !match {
				writer.WriteError(w, structs.ErrUnauthorized)
				return
			}

			tokenString := strings.Split(authorization, " ")

			token, err := jwtt.Parse(tokenString[1])
			if err != nil {
				writer.WriteError(w, structs.ErrUnauthorized)
				return
			}

			claims, ok := token.Claims.(*jwt.UserClaim)
			if !ok {
				writer.WriteError(w, structs.ErrUnauthorized)
				return
			}

			parentCtx := r.Context()
			ctx := context.WithValue(parentCtx, "token", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
