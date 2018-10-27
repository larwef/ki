package auth

import (
	"encoding/base64"
	"github.com/larwef/ki/test"
	"testing"
)

func getBasicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func TestRegisterAndAuthenticate(t *testing.T) {
	basic := NewBasic()

	hashPassword, err := HashPassword("testPassword")
	test.AssertNotError(t, err)

	user := User{
		Username:     "someUserName",
		PasswordHash: hashPassword,
		Role:         CLIENT,
	}

	err = basic.RegisterUser(user)
	test.AssertNotError(t, err)

	authUser, err := basic.Authenticate(getBasicAuth("someUserName", "testPassword"))
	test.AssertNotError(t, err)
	test.AssertEqual(t, authUser.Username, "someUserName")
	test.AssertEqual(t, authUser.PasswordHash, user.PasswordHash)
	test.AssertEqual(t, authUser.Role, CLIENT)
}

func TestRegisterUser_DuplicateUserName(t *testing.T) {
	basic := NewBasic()

	hashPassword1, err := HashPassword("testPassword1")
	test.AssertNotError(t, err)

	user1 := User{
		Username:     "someUserName",
		PasswordHash: hashPassword1,
		Role:         CLIENT,
	}

	err = basic.RegisterUser(user1)
	test.AssertNotError(t, err)

	hashPassword2, err := HashPassword("testPassword2")
	test.AssertNotError(t, err)

	user2 := User{
		Username:     "someUserName",
		PasswordHash: hashPassword2,
		Role:         ADMIN,
	}

	err = basic.RegisterUser(user2)
	test.AssertEqual(t, err, ErrUserAlreadyExists)
}

func TestAuthenticate_UserDoesntExist(t *testing.T) {
	basic := NewBasic()

	_, err := basic.Authenticate(getBasicAuth("user", "password"))
	test.AssertEqual(t, err, ErrUserDoesntExist)
}

func TestAuthenticate_PasswordInvalid(t *testing.T) {
	basic := NewBasic()

	hashPassword, err := HashPassword("testPassword")
	test.AssertNotError(t, err)

	user := User{
		Username:     "someUserName",
		PasswordHash: hashPassword,
		Role:         CLIENT,
	}

	err = basic.RegisterUser(user)
	test.AssertNotError(t, err)

	_, err = basic.Authenticate(getBasicAuth("someUserName", "wrongPassword"))
	test.AssertEqual(t, err, ErrInvalidPassword)
}
