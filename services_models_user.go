package main

import (
	"errors"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Get first user by id
func GetUserByID(uid uint) CpsAdmins {
	var u CpsAdmins
	if err := db_ksc.First(&u, uid).Error; err != nil {
		return u
	}
	u.PrepareGive()
	return u
}

func (u *CpsAdmins) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Login check func
func LoginCheck(email string, password string) (string, CpsAdmins, error) {
	var err error
	u := CpsAdmins{}
	err = db_ksc.Model(CpsAdmins{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", u, err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", u, err
	}
	token, err := GenerateToken(u.ID)
	if err != nil {
		return "", u, err
	}
	return token, u, nil

}

// Save user and first create
func SaveUser(u CpsAdmins) error {
	user := CpsAdmins{}
	resp := db_ksc.Model(CpsAdmins{}).Where("email = ? or phone = ?", u.Email, u.Phone).Find(&user)
	if resp.RowsAffected > 0 {
		err := errors.New("account has been taken")
		return err
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	err := db_ksc.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *CpsAdmins) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	//remove spaces in username
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil

}
