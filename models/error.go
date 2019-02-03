package models

import (
	"errors"
	"fmt"
)

var (
	errServer = errors.New("Internal server error")
)

// ErrUserAlreadyExist represents a "User already exists" error.
type ErrUserAlreadyExist struct {
	Email string
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("User already exists [email: %s]", err.Email)
}

// ErrUserNotExist represents a "User already exists" error.
type ErrUserNotExist struct {
	Email string
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("User does not exist [email: %s]", err.Email)
}

// ErrUserWrongEmail represents a "User's email is wrong" error.
type ErrUserWrongEmail struct {
	Email string
}

func (err ErrUserWrongEmail) Error() string {
	return fmt.Sprintf("User's email is wrong [email: %s]", err.Email)
}

// ErrUserWrongName represents a "User's name is wrong" error.
type ErrUserWrongName struct {
	Firstname string
	Lastname  string
}

func (err ErrUserWrongName) Error() string {
	return fmt.Sprintf("User's name is wrong [firstname: %s, lastname: %s]", err.Firstname, err.Lastname)
}

// ErrUserWrongPasswd represents a "User's passwd is wrong" error.
type ErrUserWrongPasswd struct {
	Passwd    string
	NewPasswd string
}

func (err ErrUserWrongPasswd) Error() string {
	return fmt.Sprintf("User's passwd is wrong [passwd: %s, new passwd: %s]", err.Passwd, err.NewPasswd)
}

// ErrUserTelephoneAlreadyExist represents a "User's telephone already exists" error.
type ErrUserTelephoneAlreadyExist struct {
	Telephone string
}

func (err ErrUserTelephoneAlreadyExist) Error() string {
	return fmt.Sprintf("User's telephone already exists [telephone: %s]", err.Telephone)
}

// ErrUserIsDisabled represents a "User is disabled" error.
type ErrUserIsDisabled struct {
	Email string
}

func (err ErrUserIsDisabled) Error() string {
	return fmt.Sprintf("User is disabled [email: %s]", err.Email)
}
