package givethemcookies

import (
	"context"
	"net/http"
	"strings"
)

type JwtVerifier interface {
	Verify(rawJWT string, dest ...interface{}) error
}

var MissingAuthorizationHeader = "MISSING_AUTHORIZATION_HEADER"

var MalformedAuthorizationHeader = "MALFORMED_AUTHORIZATION_HEADER"

var InvalidJWTToken = "INVALID_JWT_TOKEN"

func JWTMiddleware(handler http.HandlerFunc, verifier JwtVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			http.Error(w, MissingAuthorizationHeader, 401)
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, MalformedAuthorizationHeader, 401)
			return
		}

		var claims Claims
		if err := verifier.Verify(splitToken[1], &claims); err != nil {
			http.Error(w, InvalidJWTToken, 401)
			log.Error(err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		handler(w, r.WithContext(ctx))
	}
}
