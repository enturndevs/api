package jwt

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	var (
		claims      Claims
		tokenString string
		iat         time.Time
		exp         time.Time
		err         error
	)
	iat, _ = time.Parse("20060102150405", "20190201000000")
	exp, _ = time.Parse("20060102150405", "20300101000000")
	claims = Claims{
		"test": "test",
		"iat":  iat.Unix(),
		"exp":  exp.Unix(),
	}
	tokenString, err = Encode(claims, "key")
	assert.Nil(t, err)
	assert.Equal(t, tokenString, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImlhdCI6MTU0ODk3OTIwMCwidGVzdCI6InRlc3QifQ.Sh0dt_wVi57jvQgybXM3gxRNzT4yF4KMpKaU41xuha4")
}

func TestDecode(t *testing.T) {
	var (
		tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImlhdCI6MTU0ODk3OTIwMCwidGVzdCI6InRlc3QifQ.Sh0dt_wVi57jvQgybXM3gxRNzT4yF4KMpKaU41xuha4"
		claims      jwt.MapClaims
		err         error
	)
	claims, err = Decode(tokenString, "key")
	assert.Nil(t, err)
	assert.Equal(t, claims["test"], "test")
}
