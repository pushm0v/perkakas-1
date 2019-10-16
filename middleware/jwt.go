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

type JWTMiddleware struct {
	signKey []byte
	jwt     *jwt.JWT
	phttp.CustomWriter
}

func NewJWTMiddleware(hctx phttp.HttpHandlerContext, signKey []byte) *JWTMiddleware {
	return &JWTMiddleware{
		signKey: signKey,
		jwt:     jwt.NewJWT(signKey),
		CustomWriter: phttp.CustomWriter{
			C: hctx,
		},
	}
}

func (j *JWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorization := r.Header.Get("Authorization")
	match, err := regexp.MatchString("^Bearer .+", authorization)
	if err != nil || !match {
		j.WriteError(w, structs.ErrUnauthorized)
		return
	}

	tokenString := strings.Split(authorization, " ")

	token, err := j.jwt.Parse(tokenString[1])
	if err != nil {
		j.WriteError(w, structs.ErrUnauthorized)
		return
	}

	claims, ok := token.Claims.(*jwt.UserClaim)
	if !ok {
		j.WriteError(w, structs.ErrUnauthorized)
		return
	}

	parentCtx := r.Context()
	ctx := context.WithValue(parentCtx, "token", claims)

	next(w, r.WithContext(ctx))
}
