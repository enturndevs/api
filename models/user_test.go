package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUserID    = int64(1)
	testFirstname = "Stanly"
	testLastname  = "Liao"
	testEmail     = "test@gmail.com"
	testPasswd    = "test1"
	testNewPasswd = "test2"
	testTelephone = "111111111"
)

func TestCryptoPassword(t *testing.T) {
	passwdSha1 := CryptoPasswd(testPasswd)
	assert.Equal(t, len(passwdSha1), 40)
}

func TestCreateUser(t *testing.T) {
	t.Run("Create a correct user", func(t *testing.T) {
		err := CreateUser(&User{
			Firstname: testFirstname,
			Lastname:  testLastname,
			Email:     testEmail,
			Passwd:    testPasswd,
			Telephone: testTelephone,
		})
		assert.Nil(t, err)
	})
	t.Run("Create a exist user", func(t *testing.T) {
		err := CreateUser(&User{
			Firstname: testFirstname,
			Lastname:  testLastname,
			Email:     testEmail,
			Passwd:    testPasswd,
		})
		assert.Equal(t, err, ErrUserAlreadyExist{testEmail})
	})
	t.Run("Create a user without email", func(t *testing.T) {
		err := CreateUser(&User{
			Firstname: testFirstname,
			Lastname:  testLastname,
			Passwd:    testPasswd,
		})
		assert.Equal(t, err, ErrUserWrongEmail{})
	})
	t.Run("Create a user without name", func(t *testing.T) {
		err := CreateUser(&User{
			Email: testEmail,
		})
		assert.Equal(t, err, ErrUserWrongName{})
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("Get existing user by id", func(t *testing.T) {
		user, err := GetUserByID(testUserID)
		assert.Nil(t, err)
		assert.Equal(t, user.Email, testEmail)
		assert.Equal(t, user.Firstname, testFirstname)
		assert.Equal(t, user.Lastname, testLastname)
		assert.Equal(t, user.Telephone, testTelephone)
	})
	t.Run("Get nonexistent user by id", func(t *testing.T) {
		id := int64(2)
		user, err := GetUserByID(id)
		assert.Nil(t, err)
		assert.Nil(t, user)
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("Get existing user by email", func(t *testing.T) {
		user, err := GetUserByEmail(testEmail)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, testUserID)
		assert.Equal(t, user.Firstname, testFirstname)
		assert.Equal(t, user.Lastname, testLastname)
		assert.Equal(t, user.Telephone, testTelephone)
	})
	t.Run("Get nonexistent user by email", func(t *testing.T) {
		email := "test@test"
		user, err := GetUserByEmail(email)
		assert.Nil(t, err)
		assert.Nil(t, user)
	})
	t.Run("Get user by wrong email", func(t *testing.T) {
		email := "test"
		user, err := GetUserByEmail(email)
		assert.Nil(t, user)
		assert.Equal(t, err, ErrUserWrongEmail{email})
	})
}

func TestGetUserByTelephone(t *testing.T) {
	t.Run("Get existing user by telephone", func(t *testing.T) {
		user, err := GetUserByTelephone(testTelephone)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, testUserID)
		assert.Equal(t, user.Email, testEmail)
		assert.Equal(t, user.Firstname, testFirstname)
		assert.Equal(t, user.Lastname, testLastname)
	})
	t.Run("Get nonexistent user by telephone", func(t *testing.T) {
		telephone := "222222222"
		user, err := GetUserByTelephone(telephone)
		assert.Nil(t, err)
		assert.Nil(t, user)
	})
}

func TestIncrUserCashback(t *testing.T) {
	cashback, err := IncrUserCashback(testUserID, 10)
	assert.Nil(t, err)
	assert.Equal(t, cashback, uint32(10))
}

func TestDecrUserCashback(t *testing.T) {
	cashback, err := DecrUserCashback(testUserID, 5)
	assert.Nil(t, err)
	assert.Equal(t, cashback, uint32(5))
}

func TestUpdateUserPasswd(t *testing.T) {
	t.Run("Update user's correct passwd", func(t *testing.T) {
		err := UpdateUserPasswd(1, testPasswd, testNewPasswd)
		assert.Nil(t, err)
		err = UpdateUserPasswd(1, testNewPasswd, testPasswd)
		assert.Nil(t, err)
	})
	t.Run("Update user's wrong passwd", func(t *testing.T) {
		err := UpdateUserPasswd(1, testNewPasswd, testPasswd)
		assert.Equal(t, err, ErrUserWrongPasswd{Passwd: testNewPasswd})
	})
}

func TestDisableUser(t *testing.T) {
	err := DisableUser(testUserID)
	assert.Nil(t, err)
}

func TestEnableUser(t *testing.T) {
	err := EnableUser(testUserID)
	assert.Nil(t, err)
}

func TestLoginUser(t *testing.T) {
	t.Run("User login is correct", func(t *testing.T) {
		user, err := LoginUser(testEmail, testPasswd)
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})
	t.Run("User login is wrong", func(t *testing.T) {
		wrongEmail := "a@a123"
		user, err := LoginUser(wrongEmail, testPasswd)
		assert.Equal(t, err, ErrUserNotExist{wrongEmail})
		assert.Nil(t, user)
	})
}

func TestUpdateUser(t *testing.T) {
	data := map[string]interface{}{
		"firstname": "Feng",
		"lastname":  "Chen",
		"telephone": testTelephone,
		"address":   "NY, USA",
	}
	t.Run("Update user by exist telephone", func(t *testing.T) {
		err := UpdateUser(testUserID, data)
		assert.Equal(t, err, ErrUserTelephoneAlreadyExist{testTelephone})
	})
	t.Run("Update user by nonexistent telephone", func(t *testing.T) {
		data["telephone"] = "222222222"
		err := UpdateUser(testUserID, data)
		assert.Nil(t, err)
	})
}
