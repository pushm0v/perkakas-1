package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/token/jwt"
)

type JWTMiddleware struct {
	signKey           []byte
	jwt               *jwt.JWT
	unauthorizedError error
	phttp.CustomWriter
}

func NewJWTMiddleware(hctx phttp.HttpHandlerContext, unauthorizedError error, signKey []byte) *JWTMiddleware {
	return &JWTMiddleware{
		signKey:           signKey,
		jwt:               jwt.NewJWT(signKey),
		unauthorizedError: unauthorizedError,
		CustomWriter: phttp.CustomWriter{
			C: hctx,
		},
	}
}

// func NewHandlerContext()

func (j *JWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorization := r.Header.Get("Authorization")
	match, err := regexp.MatchString("^Bearer .+", authorization)
	if err != nil || !match {
		j.WriteError(w, j.unauthorizedError)
		return
	}

	tokenString := strings.Split(authorization, " ")

	token, err := j.jwt.Parse(tokenString[1])
	if err != nil {
		j.WriteError(w, j.unauthorizedError)
		return
	}

	claims, ok := token.Claims.(*jwt.UserClaim)
	if !ok {
		j.WriteError(w, j.unauthorizedError)
		return
	}

	parentCtx := r.Context()
	ctx := context.WithValue(parentCtx, "token", claims)

	next(w, r.WithContext(ctx))

	// err := common.ErrInvalidRequestEntity
	// auth := strings.TrimSpace(r.Header.Get("Authorization"))
	// if auth == "" {
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }

	// authToken := strings.Split(auth, " ")
	// if len(authToken) != 2 {
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }

	// tokenType := authToken[0]
	// tokenValue := authToken[1]

	// if tokenType != "Bearer" {
	// 	err = common.ErrInvalidToken
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }

	// // logic for support old token (the long one). Should be refactor this block after ramadan
	// splittedToken := strings.Split(decryptedToken, ".")
	// content, err := jwt.DecodeSegment(splittedToken[1])
	// if err != nil {
	// 	metric.CountJwtTokenInvalid("Fail Decode Segment Old Token")
	// 	err = common.ErrInvalidToken
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }

	// var newClaim model.NewKitabisaClaims
	// err = json.Unmarshal(content, &newClaim)
	// if err != nil {
	// 	metric.CountJwtTokenInvalid("Fail Unmarshal New Claims")
	// 	err = common.ErrInvalidToken
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }

	// var token *jwt.Token
	// if newClaim.UserID != 0 {
	// 	token, err = jwt.ParseWithClaims(decryptedToken, &model.NewKitabisaClaims{}, func(t *jwt.Token) (interface{}, error) {
	// 		return []byte(m.JWTSignKey), nil
	// 	})
	// 	if err != nil {
	// 		if m.IsTokenExpired(err) {
	// 			metric.CountJwtTokenInvalid("Token Expired")
	// 		} else {
	// 			metric.CountJwtTokenInvalid("Fail Parse New Claims")
	// 		}

	// 		err = common.ErrInvalidToken
	// 		statusError := handler.StatusError{
	// 			Code: http.StatusUnauthorized,
	// 			Err: err,
	// 		}
	// 		handler.WriteErrorResponse(w, statusError)
	// 		return
	// 	}
	// } else {
	// 	token, err = jwt.ParseWithClaims(decryptedToken, &model.KitabisaClaims{}, func(t *jwt.Token) (interface{}, error) {
	// 		return []byte(m.JWTSignKey), nil
	// 	})
	// 	if err != nil {
	// 		if m.IsTokenExpired(err) {
	// 			metric.CountJwtTokenInvalid("Token Expired")
	// 		} else {
	// 			metric.CountJwtTokenInvalid("Fail Parse Old Claims")
	// 		}
	// 		err = common.ErrInvalidToken
	// 		statusError := handler.StatusError{
	// 			Code: http.StatusUnauthorized,
	// 			Err: err,
	// 		}
	// 		handler.WriteErrorResponse(w, statusError)
	// 		return
	// 	}
	// }

	// ctx := r.Context()
	// if claims, ok := token.Claims.(*model.NewKitabisaClaims); ok && token.Valid {
	// 	userID := claims.UserID
	// 	_, err := m.TokenService.GetToken(userID)
	// 	if err != nil {
	// 		if err == redis.ErrNil {
	// 			err = common.ErrSessionExpired
	// 		}
	// 		statusError := handler.StatusError{
	// 			Code: http.StatusUnauthorized,
	// 			Err: err,
	// 		}
	// 		handler.WriteErrorResponse(w, statusError)
	// 		return
	// 	}

	// 	ctx = context.WithValue(ctx, "UserInfo", claims)
	// } else if claims, ok := token.Claims.(*model.KitabisaClaims); ok && token.Valid {
	// 	userID := claims.User.ID
	// 	_, err := m.TokenService.GetToken(userID)
	// 	if err != nil {
	// 		if err == redis.ErrNil {
	// 			err = common.ErrSessionExpired
	// 		}
	// 		statusError := handler.StatusError{
	// 			Code: http.StatusUnauthorized,
	// 			Err: err,
	// 		}
	// 		handler.WriteErrorResponse(w, statusError)
	// 		return
	// 	}

	// 	ctx = context.WithValue(ctx, "UserInfo", claims.User)
	// } else {
	// 	statusError := handler.StatusError{
	// 		Code: http.StatusUnauthorized,
	// 		Err: err,
	// 	}
	// 	handler.WriteErrorResponse(w, statusError)
	// 	return
	// }
	// // end of logic for support old token

	// next(w, r.WithContext(ctx))
}

// func (m *Middleware) IsTokenExpired(err error) bool {
// 	switch err.(type) {
// 	case *jwt.ValidationError:
// 		vErr := err.(*jwt.ValidationError)
// 		switch vErr.Errors {
// 		case jwt.ValidationErrorExpired:
// 			return true
// 		default:
// 			return false
// 		}

// 	default:
// 		return false
// 	}
// }
