package util

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Store - secure cookie store
var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

func init() {
	Store.Options.HttpOnly = true
	Store.MaxAge(3600 * 24) // max age is 24 hours of log tailing
}

// Login - tries to authenticate the user via PAM
func Login(r *http.Request) (bool, string, error) {
	username, password := r.FormValue("username"), r.FormValue("password")

	if !IsWhitelisted(username) {
		return false, "", fmt.Errorf("username %s denied due to ACL", username)
	}

	isValid := PamAauthenticate(username, password)
	if isValid == 1 {
		return true, username, nil
	}
	return false, "", errors.New("Invalid Login credentials")
}

// GenerateSecureKey - Key for CSRF Tokens
func GenerateSecureKey() string {
	// Inspired from gorilla/securecookie
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}
