package givethemcookies

import (
	"context"
	"net/http"
	"strings"
)

// JWTVerify verifies JWTs and populates "dest" with claims
type JwtVerifier interface {
	Verify(rawJWT string, dest ...interface{}) error
}

// ClaimFetcher is a function which returns a pointer for the claims to be populated against, e.g. givethemcookies.Claims
type ClaimFetcher = func() interface{}

// MissingAuthorizationHeader is returned to the http client when there's a missing authorization header in the request
var MissingAuthorizationHeader = "MISSING_AUTHORIZATION_HEADER"

// MalformedAuthorizationHeader is returned to the http client when there's a malformed authorization header in the request
var MalformedAuthorizationHeader = "MALFORMED_AUTHORIZATION_HEADER"

// InvalidJWTToken is returned to the http client when an invalid token is supplied in the request
var InvalidJWTToken = "INVALID_JWT_TOKEN"

// JWTMiddleware authorizes an populates claims from claimFetcher and sets it to context
func JWTMiddleware(handler http.HandlerFunc, verifier JwtVerifier, claimFetcher ClaimFetcher) http.HandlerFunc {
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

		claim := claimFetcher()
		if err := verifier.Verify(splitToken[1], claim); err != nil {
			http.Error(w, InvalidJWTToken, 401)
			log.Error(err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claim)
		handler(w, r.WithContext(ctx))
	}
}
