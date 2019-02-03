package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/enturndevs/enturn-service/modules/jwt"
	"github.com/spf13/viper"
	"github.com/stanlyliao/logger"

	"github.com/enturndevs/enturn-service/models"

	"github.com/gin-gonic/gin"
)

const (
	// ErrUserLoginData represents user's data is wrong
	ErrUserLoginData = "Error Data"
)

// ResponseLogin represents login response
type ResponseLogin struct {
	Response
	Data map[string]interface{}
}

func initUserRoutes() {
	v1 := r.Group("/v1")
	user := v1.Group("/user")
	user.POST("/signup", userSignup)
	user.POST("/login", userLogin)
}

func userSignup(c *gin.Context) {
	reqData, _ := ioutil.ReadAll(c.Request.Body)

	var data Request
	if err := json.Unmarshal(reqData, &data); err != nil {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}

	user := models.User{}

	email, ok := data["email"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}
	user.Email = email.(string)

	firstname, ok := data["firstname"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}
	user.Firstname = firstname.(string)

	lastname, ok := data["lastname"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}
	user.Lastname = lastname.(string)

	passwd, ok := data["passwd"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}
	user.Passwd = passwd.(string)

	telephone, ok := data["telephone"]
	if ok {
		user.Telephone = telephone.(string)
	}

	address, ok := data["address"]
	if ok {
		user.Address = address.(string)
	}

	user.IsActive = true

	if err := models.CreateUser(&user); err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := Response{
		Success: true,
	}
	c.JSON(http.StatusOK, resp)
}

func userLogin(c *gin.Context) {
	reqData, _ := ioutil.ReadAll(c.Request.Body)

	var data Request
	if err := json.Unmarshal(reqData, &data); err != nil {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}

	email, ok := data["email"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}

	passwd, ok := data["passwd"]
	if !ok {
		responseError(c, http.StatusBadRequest, ErrUserLoginData)
		return
	}

	user, err := models.LoginUser(email.(string), passwd.(string))
	if err != nil {
		logger.Error(err)
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	claims := jwt.Claims{
		"email": user.Email,
	}
	token, err := jwt.Encode(claims, viper.GetString("server.token"))
	if err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := ResponseLogin{
		Response: Response{
			Success: true,
		},
		Data: map[string]interface{}{
			"token": token,
		},
	}
	c.JSON(http.StatusOK, resp)
}
