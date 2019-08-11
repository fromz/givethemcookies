package givethemcookies

import (
	"errors"
	"fmt"
	"github.com/fromz/go-auth"
	"os"
)

func GetJWTVerifierFromENVVars() (*go_auth.JwtVerifier, error) {
	// configure the jwt verifier
	JWKSSource, found := os.LookupEnv("JWKS_SOURCE")
	if !found {
		return nil, errors.New("expected environment variable JWKS_SOURCE to be set")
	}
	JWKSLocation, found := os.LookupEnv("JWKS_LOCATION")
	if !found {
		return nil, errors.New("expected environment variable JWKS_LOCATION to be set")
	}

	JWKSKeyID, found := os.LookupEnv("JWKS_KEY_ID")
	if !found {
		return nil, errors.New("expected environment variable JWKS_KEY_ID to be set")
	}

	return getJWTVerifier(JWKSSource, JWKSLocation, JWKSKeyID)
}

func getJWTVerifier(source, location, JWKSKeyID string) (*go_auth.JwtVerifier, error) {
	switch source {
	case "file":
		fs, err := go_auth.JWKSFileSource(location)
		if err != nil {
			return nil, err
		}
		JWTVerifier := go_auth.NewJwtVerifier(go_auth.JWKSClient(fs), JWKSKeyID)
		return &JWTVerifier, nil
	case "http":
		fs := go_auth.JWKSWebSource(location)
		JWTVerifier := go_auth.NewJwtVerifier(go_auth.JWKSClient(fs), JWKSKeyID)
		return &JWTVerifier, nil
	default:
		return nil, fmt.Errorf("JWKS_SOURCE can be either file or http, received: %s", source)
	}
}
