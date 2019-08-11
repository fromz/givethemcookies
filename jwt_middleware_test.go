package givethemcookies_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/fromz/givethemcookies"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// GetTestHandler returns a http.HandlerFunc for testing http middleware
func GetTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("ACCESSED_ENDPOINT"))
	}
	return http.HandlerFunc(fn)
}

type successJwtVerifier struct {
}

func (s successJwtVerifier) Verify(rawJWT string, dests ...interface{}) error {
	for _, dest := range dests {
		switch v := dest.(type) {
		case *givethemcookies.Claims:
			v.Subject = "Kieran"
		}

	}
	return nil
}

type failureJwtVerifier struct {
}

func (s failureJwtVerifier) Verify(rawJWT string, dest ...interface{}) error {
	return errors.New("JWT Verification Failed")
}

func prepareTest(t *testing.T, handlerFunc http.HandlerFunc, verifier givethemcookies.JwtVerifier, jwtHeader *string) (string, int) {
	givethemcookies.SetLogger(logrus.New())
	ts := httptest.NewServer(givethemcookies.JWTMiddleware(handlerFunc, verifier))
	defer ts.Close()

	var u bytes.Buffer
	u.WriteString(string(ts.URL))
	u.WriteString("/")

	req, err := http.NewRequest("GET", u.String(), nil)
	if jwtHeader != nil {
		req.Header.Add("Authorization", *jwtHeader)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	return string(b), res.StatusCode
}

func TestJWTMiddleware_MissingAuthorizationBearer(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), failureJwtVerifier{}, nil)
	s = strings.TrimSpace(s)
	if s != givethemcookies.MissingAuthorizationHeader {
		t.Errorf("Expected `%s` got `%s`", givethemcookies.MissingAuthorizationHeader, s)
	}
	if code != 401 {
		t.Errorf("Expected %d got %d", 401, code)
	}
}

func TestJWTMiddleware_EmptyAuthorizationBearer(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), failureJwtVerifier{}, stringToPointerString(""))
	s = strings.TrimSpace(s)
	if s != givethemcookies.MissingAuthorizationHeader {
		t.Errorf("Expected `%s` got `%s`", givethemcookies.MissingAuthorizationHeader, s)
	}
	if code != 401 {
		t.Errorf("Expected %d got %d", 401, code)
	}
}

func TestJWTMiddleware_MalformedAuthorizationBearer_WithSpace(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), failureJwtVerifier{}, stringToPointerString("Bearer "))
	s = strings.TrimSpace(s)
	if s != givethemcookies.MalformedAuthorizationHeader {
		t.Errorf("Expected `%s` got `%s`", givethemcookies.MalformedAuthorizationHeader, s)
	}
	if code != 401 {
		t.Errorf("Expected %d got %d", 401, code)
	}
}

func TestJWTMiddleware_MalformedAuthorizationBearer_WithoutSpace(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), failureJwtVerifier{}, stringToPointerString("Bearer"))
	s = strings.TrimSpace(s)
	if s != givethemcookies.MalformedAuthorizationHeader {
		t.Errorf("Expected `%s` got `%s`", givethemcookies.MalformedAuthorizationHeader, s)
	}
	if code != 401 {
		t.Errorf("Expected %d got %d", 401, code)
	}
}

func TestJWTMiddleware_InvalidJWT(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), failureJwtVerifier{}, stringToPointerString("Bearer INVALID.JWT.TOKEN"))
	s = strings.TrimSpace(s)
	if s != givethemcookies.InvalidJWTToken {
		t.Errorf("Expected `%s` got `%s`", givethemcookies.InvalidJWTToken, s)
	}
	if code != 401 {
		t.Errorf("Expected %d got %d", 401, code)
	}
}

func TestJWTMiddleware_ValidJWT(t *testing.T) {
	s, code := prepareTest(t, GetTestHandler(), successJwtVerifier{}, stringToPointerString("Bearer VALID.JWT.TOKEN"))
	s = strings.TrimSpace(s)
	if s != "ACCESSED_ENDPOINT" {
		t.Errorf("Expected `%s` got `%s`", "ACCESSED_ENDPOINT", s)
	}
	if code != 200 {
		t.Errorf("Expected %d got %d", 200, code)
	}
}

// GetClaimsToJsonHandler returns a http.HandlerFunc for testing http middleware
func GetClaimsToJsonHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		user := req.Context().Value("user")
		switch v := user.(type) {
		case givethemcookies.Claims:
			b, err := json.Marshal(v)
			if err != nil {
				panic(err.Error())
			}
			rw.Write(b)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func TestJWTMiddleware_ValidJWT_ClaimsAreSet(t *testing.T) {
	s, code := prepareTest(t, GetClaimsToJsonHandler(), successJwtVerifier{}, stringToPointerString("Bearer VALID.JWT.TOKEN"))
	var claims givethemcookies.Claims
	if err := json.Unmarshal([]byte(s), &claims); err != nil {
		t.Error(err)
	}
	if claims.Subject != "Kieran" {
		t.Errorf("expected subject to be %s got %s", "Kieran", claims.Subject)
	}
	if code != 200 {
		t.Errorf("Expected %d got %d", 200, code)
	}
}

func stringToPointerString(s string) *string {
	return &s
}
