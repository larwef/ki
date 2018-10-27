package auth

import (
	"errors"
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

// ErrUserDoesntExist is used when trying to authenticate a user which doesnt exist in the pool.
var ErrUserDoesntExist = errors.New("user doesnt exists")

// ErrInvalidPassword is used when trying to authenticate a user using an incorrect password
var ErrInvalidPassword = errors.New("invalid password")

// User represents a user object.
type User struct {
	Username     string
	PasswordHash string
	Role         Role
}
