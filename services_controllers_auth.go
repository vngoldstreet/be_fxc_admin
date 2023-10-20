package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get current user by token auth request
func CurrentUser(c *gin.Context) {
	user_id, err := ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "code": http.StatusOK, "data": u})
}

// Login func
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := CpsAdmins{}

	u.Email = input.Email
	u.Password = input.Password
	token, u, err := LoginCheck(u.Email, u.Password)
	u.Password = ""
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
}

// Register func
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := CpsAdmins{}

	u.Name = input.Name
	u.Email = input.Email
	u.Password = input.Password
	err := SaveUser(u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failure",
			"message": "Account already exists",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Registration success",
		})
	}
}
