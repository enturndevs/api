package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testToken     string
	testFirstname = "Stanly"
	testLastname  = "Liao"
	testEmail     = "test2@gmail.com"
	testPasswd    = "test1"
	testNewPasswd = "test2"
	testTelephone = "111111111"
)

func TestUserSignup(t *testing.T) {
	var resp ResponseLogin
	t.Run("User signup without name", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte(fmt.Sprintf(`{"email":"%s","passwd":"%s"}`, testEmail, testPasswd))
		req, err := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(data))
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, resp.ErrorMsg, ErrUserLoginData)
	})
	t.Run("User signup without email", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte(fmt.Sprintf(`{"firstname":"%s","lastname":"%s","passwd":"%s"}`, testFirstname, testLastname, testPasswd))
		req, err := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(data))
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.False(t, resp.Success)
		assert.Equal(t, resp.ErrorMsg, ErrUserLoginData)
	})
	t.Run("User signup", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte(fmt.Sprintf(`{"email":"%s","firstname":"%s","lastname":"%s","passwd":"%s"}`, testEmail, testFirstname, testLastname, testPasswd))
		req, err := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(data))
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
	})
}

func TestUserLogin(t *testing.T) {
	var resp ResponseLogin
	t.Run("User login by correct data", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte(fmt.Sprintf(`{"email":"%s","passwd":"%s"}`, testEmail, testPasswd))
		req, err := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(data))
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.Data["token"])

		testToken = resp.Data["token"].(string)
	})
}
