package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	BASIC_SCHEMA  string = "Basic "
	BEARER_SCHEMA string = "Bearer "
)

// FromRequest extracts the auth token from req.
func FromRequest(req *http.Request) (string, error) {
	// Grab the raw Authoirzation header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	// Confirm the request is sending Basic Authentication credentials.
	if !strings.HasPrefix(authHeader, BASIC_SCHEMA) && !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
		return "", errors.New("Authorization requires Basic/Bearer scheme")
	}

	// Get the token from the request header
	// The first six characters are skipped - e.g. "Basic ".
	if strings.HasPrefix(authHeader, BASIC_SCHEMA) {
		str, err := base64.StdEncoding.DecodeString(authHeader[len(BASIC_SCHEMA):])
		if err != nil {
			return "", errors.New("Base64 encoding issue")
		}
		creds := strings.Split(string(str), ":")
		return creds[0], nil
	} else {
		return authHeader[len(BEARER_SCHEMA):], nil
	}
}

func SetRequestHeader(req *http.Request, token string){
	req.Header.Set("Authorization", "Bearer TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
}
