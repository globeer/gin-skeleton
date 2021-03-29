package controllers

import (
	"net/http"

	"bower.co.kr/c4bapi/models"
	"github.com/gin-gonic/gin"
)

// UserController is the user controller
type User1Controller struct{}

// Signup struct
type Signup struct {
	UserId  string `form:"userId" json:"user_id" binding:"required"`
	UserNm  string `form:"userNm" json:"user_nm" binding:"required"`
	UserPw  string `form:"userPw" json:"user_pw" binding:"required,min=6,max=20"`
	UserPw2 string `form:"userPw2" json:"user_pw2" binding:"required"`
}

// GetUser gets the user info
func (ctrl *User1Controller) GetUser(c *gin.Context) {
	var user models.User1

	userId := c.Param("userId")

	if err := user.GetFirstByID(userId); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// SignupForm shows the signup form
func (ctrl *User1Controller) SignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

// LoginForm shows the login form
func (ctrl *User1Controller) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Signup a new user
func (ctrl *User1Controller) Signup(c *gin.Context) {
	var form Signup
	if err := c.ShouldBind(&form); err == nil {

		if form.UserPw != form.UserPw2 {
			c.JSON(http.StatusOK, gin.H{"error": "Password does not match with conform password"})
			return
		}

		var user models.User1

		user.UserNm = form.UserNm
		user.UserId = form.UserId
		user.UserPw = form.UserPw

		if err := user.Signup(); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
