package auth

import (
	"github.com/larwef/ki/test"
	"testing"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	pool := NewUserPool()

	hashPassword, err := HashPassword("testPassword")
	test.AssertNotError(t, err)

	user := User{
		Username:     "someUserName",
		PasswordHash: hashPassword,
		Role:         CLIENT,
	}

	err = pool.RegisterUser(user)
	test.AssertNotError(t, err)

	authUser, err := pool.Authenticate("someUserName", "testPassword")
	test.AssertNotError(t, err)
	test.AssertEqual(t, authUser.Username, "someUserName")
	test.AssertEqual(t, authUser.PasswordHash, user.PasswordHash)
	test.AssertEqual(t, authUser.Role, CLIENT)
}

func TestRegisterUser_DuplicateUserName(t *testing.T) {
	pool := NewUserPool()

	hashPassword1, err := HashPassword("testPassword1")
	test.AssertNotError(t, err)

	user1 := User{
		Username:     "someUserName",
		PasswordHash: hashPassword1,
		Role:         CLIENT,
	}

	err = pool.RegisterUser(user1)
	test.AssertNotError(t, err)

	hashPassword2, err := HashPassword("testPassword2")
	test.AssertNotError(t, err)

	user2 := User{
		Username:     "someUserName",
		PasswordHash: hashPassword2,
		Role:         ADMIN,
	}

	err = pool.RegisterUser(user2)
	test.AssertEqual(t, err, ErrUserAlreadyExists)
}

func TestAuthenticate_UserDoesntExist(t *testing.T) {
	pool := NewUserPool()

	_, err := pool.Authenticate("user", "password")
	test.AssertEqual(t, err, ErrUserDoesntExist)
}

func TestAuthenticate_PasswordInvalid(t *testing.T) {
	pool := NewUserPool()

	hashPassword, err := HashPassword("testPassword")
	test.AssertNotError(t, err)

	user := User{
		Username:     "someUserName",
		PasswordHash: hashPassword,
		Role:         CLIENT,
	}

	err = pool.RegisterUser(user)
	test.AssertNotError(t, err)

	_, err = pool.Authenticate("someUserName", "wrongPassword")
	test.AssertEqual(t, err, ErrInvalidPassword)
}
