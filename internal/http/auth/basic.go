package auth

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// Basic is used to provide basic authentication.
type Basic struct {
	userPool map[string]User
}

// NewBasic returns a new Basic object.
func NewBasic() *Basic {
	return &Basic{
		userPool: make(map[string]User),
	}
}

// Authenticate chechs that an authentication header can be successfully mapped to an existing user. Returns the appropriate User
// and nil error. Returns and error and an empty user object if not successful.
func (b *Basic) Authenticate(authHeader string) (User, error) {
	username, password, ok := parseBasicAuth(authHeader)
	if !ok {
		return User{}, errors.New("error parsing auth header")
	}

	user, exists := b.userPool[username]

	if !exists {
		return User{}, ErrUserDoesntExist
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return User{}, ErrInvalidPassword
	}

	return user, nil
}

// RegisterUser puts a user in the user pool. Will return an error if not successful.
func (b *Basic) RegisterUser(user User) error {
	if _, exists := b.userPool[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	b.userPool[user.Username] = user
	return nil
}

// Stolen from net/http.request. Parses a basic auth request header.
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// HashPassword is a helper function for calling bcrypt.GenerateFromPassword
func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}
