package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Role defines available roles
type Role int

const (
	// ADMIN role has access to both read and write operations for all resources.
	ADMIN Role = 1 << iota
	// CLIENT role has access to read operations only.
	CLIENT
)

// ErrUserAlreadyExists is used when trying to register a username which already exists in the pool.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserDoesntExist is used when trying to Authenticate a user which doesnt exist in the pool.
var ErrUserDoesntExist = errors.New("user doesnt exists")

// ErrInvalidPassword is used when trying to Authenticate a user using an incorrect password
var ErrInvalidPassword = errors.New("invalid password")

// User represents a user object.
type User struct {
	Username     string
	PasswordHash string
	Role         Role
}

// UserPool is a collection of users.
type UserPool struct {
	users map[string]User
}

func NewUserPool() *UserPool {
	return &UserPool{
		users: make(map[string]User),
	}
}

// RegisterUser adds a user to the pool. Returns error if the user couldn be added.
func (up *UserPool) RegisterUser(user User) error {
	if _, exists := up.users[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	up.users[user.Username] = user
	return nil
}

// Authenticate checks if a user exists in the pool, verifies password match and returns the matching User object.
func (up *UserPool) Authenticate(username string, password string) (User, error) {
	user, exists := up.users[username]

	if !exists {
		return User{}, ErrUserDoesntExist
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return User{}, ErrInvalidPassword
	}

	return user, nil
}

// HashPassword is a helperfunction for calling bcrypt.GenerateFromPassword
func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}
