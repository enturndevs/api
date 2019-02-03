package models

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/stanlyliao/logger"
)

// User represents the object of users
type User struct {
	ID            int64  `xorm:"PK AUTOINCR 'id'" json:"id"`
	Email         string `xorm:"UNIQUE VARCHAR(50) NOT NULL" json:"email"`
	Firstname     string `xorm:"VARCHAR(20) NOT NULL" json:"firstname"`
	Lastname      string `xorm:"VARCHAR(20) NOT NULL" json:"lastname"`
	Passwd        string `xorm:"CHAR(40) NOT NULL" json:"passwd"`
	Telephone     string `xorm:"UNIQUE VARCHAR(15)" json:"telephone"`
	Address       string `xorm:"VARCHAR(100)" json:"address"`
	Cashback      uint32 `xorm:"DEFAULT 0" json:"-"`
	IsActive      bool   `xorm:"INDEX DEFAULT false" json:"is_active"`
	LastLoginUnix int64  `xorm:"INDEX" json:"last_login_unix"`
	CreatedUnix   int64  `xorm:"INDEX CREATED" json:"-"`
	UpdatedUnix   int64  `xorm:"INDEX UPDATED" json:"-"`
}

// JSON return json format user
func (u *User) JSON() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", errServer
	}

	return string(b), nil
}

func isEmailExist(e Engine, email string) (bool, error) {
	if len(email) == 0 {
		return false, nil
	}
	return e.
		Where("email=?", email).
		Get(new(User))
}

// CryptoPasswd return crypto sha1 passwd
func CryptoPasswd(passwd string) string {
	h := sha1.New()
	h.Write([]byte(passwd))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func getUserByEmailFromRedis(email string) (*User, error) {
	key := "user:" + email
	userJSON, err := r.Get(key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, ErrUserNotExist{email}
		}
		logger.Error(err)
		return nil, errServer
	}

	var user User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		logger.Error(err)
		return nil, errServer
	}

	return &user, nil
}

func getUserByIDFromRedis(id int64) (*User, error) {
	key := fmt.Sprintf("user:%d", id)
	email, err := r.Get(key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, ErrUserNotExist{email}
		}
		logger.Error(err)
		return nil, errServer
	}

	return getUserByEmailFromRedis(email)
}

func setUserToRedis(user *User) error {
	userJSON, err := user.JSON()
	if err != nil {
		logger.Error(err)
		return errServer
	}

	key := "user:" + user.Email
	if err := r.Set(key, userJSON, 0).Err(); err != nil {
		logger.Error(err)
		return errServer
	}

	key = fmt.Sprintf("user:%d", user.ID)
	if err := r.Set(key, user.Email, 0).Err(); err != nil {
		logger.Error(err)
		return errServer
	}
	return nil
}

func delUserToRedis(user *User) error {
	key := "user:" + user.Email
	if err := r.Del(key).Err(); err != nil {
		logger.Error(err)
		return errServer
	}

	key = fmt.Sprintf("user:%d", user.ID)
	if err := r.Del(key).Err(); err != nil {
		logger.Error(err)
		return errServer
	}
	return nil
}

// CreateUser creates record of a new user
func CreateUser(user *User) (err error) {
	sess := x.NewSession()
	defer sess.Close()

	if err = sess.Begin(); err != nil {
		logger.Error(err)
		return errServer
	}

	if err = checkmail.ValidateFormat(user.Email); err != nil {
		return ErrUserWrongEmail{user.Email}
	}

	if len(user.Firstname) == 0 || len(user.Lastname) == 0 {
		return ErrUserWrongName{user.Firstname, user.Lastname}
	}

	if len(user.Passwd) == 0 {
		return ErrUserWrongPasswd{Passwd: user.Passwd}
	}

	user.Email = strings.ToLower(user.Email)
	isExist, err := isEmailExist(sess, user.Email)
	if err != nil {
		return err
	} else if isExist {
		return ErrUserAlreadyExist{user.Email}
	}

	user.Passwd = CryptoPasswd(user.Passwd)

	if _, err := sess.Insert(user); err != nil {
		logger.Error(err)
		return errServer
	}

	if err := setUserToRedis(user); err != nil {
		delUserToRedis(user)
		return errServer
	}

	if err := sess.Commit(); err != nil {
		delUserToRedis(user)
		logger.Error(err)
		return errServer
	}

	return nil
}

func getUserByID(e Engine, id int64) (*User, error) {
	user := new(User)
	has, err := e.ID(id).Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return user, nil
}

// GetUserByID return the user object by given ID if exists
func GetUserByID(id int64) (*User, error) {
	return getUserByID(x, id)
}

func getUserByEmail(e Engine, email string) (*User, error) {
	user := new(User)
	has, err := e.
		Where("email=?", email).
		Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return user, nil
}

// GetUserByEmail return the user object by given Email if exists
func GetUserByEmail(email string) (*User, error) {
	if err := checkmail.ValidateFormat(email); err != nil {
		return nil, ErrUserWrongEmail{email}
	}
	return getUserByEmail(x, email)
}

func getUserByTelephone(e Engine, telephone string) (*User, error) {
	user := new(User)
	has, err := e.
		Where("telephone=?", telephone).
		Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return user, nil
}

// GetUserByTelephone return the user object by given telephone if exists
func GetUserByTelephone(telephone string) (*User, error) {
	return getUserByTelephone(x, telephone)
}

// IncrUserCashback increase user's cashback and return
func IncrUserCashback(id int64, cashback uint32) (uint32, error) {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	if _, err := sess.
		ID(id).
		Incr("cashback", cashback).
		Update(new(User)); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	user, err := getUserByID(sess, id)
	if err != nil {
		return 0, err
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	return user.Cashback, nil
}

// DecrUserCashback decrease user's cashback and return
func DecrUserCashback(id int64, cashback uint32) (uint32, error) {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	if _, err := sess.
		ID(id).
		Decr("cashback", cashback).
		Update(new(User)); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	user, err := getUserByID(sess, id)
	if err != nil {
		return 0, err
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return 0, errServer
	}

	return user.Cashback, nil
}

// UpdateUserPasswd update user's passwd
func UpdateUserPasswd(id int64, oldPasswd string, newPasswd string) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return errServer
	}

	if len(oldPasswd) == 0 {
		return ErrUserWrongPasswd{Passwd: oldPasswd}
	}
	if len(newPasswd) == 0 {
		return ErrUserWrongPasswd{NewPasswd: newPasswd}
	}

	user, err := getUserByIDFromRedis(id)
	if err != nil {
		return err
	}

	if user.Passwd != CryptoPasswd(oldPasswd) {
		return ErrUserWrongPasswd{Passwd: oldPasswd}
	}

	user.Passwd = CryptoPasswd(newPasswd)

	if _, err := sess.
		ID(id).
		Cols("passwd").
		Update(user); err != nil {
		logger.Error(err)
		return errServer
	}

	if err := setUserToRedis(user); err != nil {
		logger.Error(err)
		return errServer
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return errServer
	}

	return nil
}

// UpdateUser updates user's data
func UpdateUser(id int64, data map[string]interface{}) error {
	user, err := getUserByIDFromRedis(id)
	if err != nil {
		logger.Error(err)
		return errServer
	}

	cols := []string{}

	firstname, ok := data["firstname"]
	if ok {
		user.Firstname = firstname.(string)
		cols = append(cols, "firstname")
	}

	lastname, ok := data["lastname"]
	if ok {
		user.Lastname = lastname.(string)
		cols = append(cols, "lastname")
	}

	telephone, ok := data["telephone"]
	if ok {
		user.Telephone = telephone.(string)
		cols = append(cols, "telephone")
	}
	u, err := GetUserByTelephone(user.Telephone)
	if err != nil {
		return err
	}
	if u != nil {
		return ErrUserTelephoneAlreadyExist{user.Telephone}
	}

	address, ok := data["address"]
	if ok {
		user.Address = address.(string)
		cols = append(cols, "address")
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return errServer
	}

	if _, err := sess.
		ID(id).
		Cols(cols...).
		Update(user); err != nil {
		logger.Error(err)
		return errServer
	}

	if err := setUserToRedis(user); err != nil {
		logger.Error(err)
		return errServer
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return errServer
	}

	return nil
}

// EnableUser enable user's status
func EnableUser(id int64) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return errServer
	}

	if _, err := sess.
		ID(id).
		Cols("is_active").
		Update(User{IsActive: true}); err != nil {
		logger.Error(err)
		return errServer
	}

	user, err := getUserByIDFromRedis(id)
	if err != nil {
		logger.Error(err)
		return errServer
	}

	user.IsActive = true
	if err != setUserToRedis(user) {
		return errServer
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return errServer
	}

	return nil
}

// DisableUser disable user's status
func DisableUser(id int64) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logger.Error(err)
		return errServer
	}

	if _, err := sess.
		ID(id).
		Cols("is_active").
		Update(&User{IsActive: false}); err != nil {
		logger.Error(err)
		return errServer
	}

	user, err := getUserByIDFromRedis(id)
	if err != nil {
		logger.Error(err)
		return errServer
	}

	user.IsActive = false
	if err != setUserToRedis(user) {
		return errServer
	}

	if err := sess.Commit(); err != nil {
		logger.Error(err)
		return errServer
	}

	return nil
}

// LoginUser user login
func LoginUser(email string, passwd string) (*User, error) {
	user, err := getUserByEmailFromRedis(email)
	if err != nil {
		return nil, err
	}

	if user.Passwd != CryptoPasswd(passwd) {
		return nil, ErrUserWrongPasswd{Passwd: passwd}
	}

	if !user.IsActive {
		return nil, ErrUserIsDisabled{email}
	}

	sess := x.NewSession()
	defer sess.Close()

	lastLoginUnix := time.Now().UTC().Unix()
	if _, err := sess.
		ID(user.ID).
		Cols("last_login_unix").
		Update(&User{LastLoginUnix: lastLoginUnix}); err != nil {
		logger.Error(err)
		return nil, errServer
	}

	return user, nil
}
