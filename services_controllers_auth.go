package main

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func resetPassword(c *gin.Context) {
	var input CpsUsers
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db_ksc.Begin()

	user := CpsUsers{}
	if err := tx.Model(CpsUsers{}).Where("email = ?", input.Email).First(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stringPassword := generateUniqueCode()
	hashedPassword, _ := HashPassword(stringPassword)
	user.Password = hashedPassword

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	SendEmailForResetPassword(user.Email, user.Email, stringPassword)
	c.JSON(http.StatusOK, gin.H{
		"status":       "Success",
		"new_password": stringPassword,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateUniqueCode() string {
	inputInt := time.Now().Unix()
	h := fnv.New32a()
	byteData := make([]byte, 4)
	byteData[0] = byte(inputInt)
	byteData[1] = byte(inputInt >> 8)
	byteData[2] = byte(inputInt >> 16)
	byteData[3] = byte(inputInt >> 24)

	h.Write(byteData)
	hashValue := h.Sum32()
	uniqueCode := fmt.Sprintf("%07X", hashValue)
	responseCode := fmt.Sprintf("%s%d@", uniqueCode, inputInt)
	return responseCode
}
