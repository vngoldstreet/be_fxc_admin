package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
)

func dbMigrations() {
	db_ksc.AutoMigrate(CpsAdmins{})
}

func CheckTokenValid(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func updateWalletForJoinContest(uid uint, cur_balance float64, new_balance float64, change float64, type_id int, contest_id string) (CpsWallets, error) {
	newTrans := CpsTransactions{}
	newTrans.TypeID = type_id
	newTrans.CustomerID = uid
	newTrans.CBalance = cur_balance
	newTrans.Amount = change
	newTrans.NBalance = new_balance
	newTrans.StatusID = 1 //Processing
	newTrans.ContestID = contest_id
	err2 := db_ksc.Create(&newTrans).Error
	if err2 != nil {
		return CpsWallets{}, err2
	}
	newWallet := CpsWallets{}
	if err := db_ksc.Model(newWallet).Where("customer_id = ?", uid).Update("balance", new_balance).Find(&newWallet).Error; err != nil {
		return CpsWallets{}, err
	}
	if err := db_ksc.Save(&newWallet).Error; err != nil {
		return CpsWallets{}, err
	}
	// newTrans.StatusID = 2

	// if err := db_ksc.Save(&newTrans).Error; err != nil {
	// 	return newWallet, err
	// }
	return newWallet, nil
}

func CheckTransactiontype(type_id int) string {
	switch type_id {
	case 1:
		return "Deposit"
	case 2:
		return "Withdraw"
	case 3:
		return "Promo"
	case 4:
		return "Join a contest"
	case 5:
		return "Earn from contest"
	default:
		return "unknow"
	}
}

func NumberToString(n int, sep rune) string {
	s := strconv.Itoa(n)
	startOffset := 0
	var buff bytes.Buffer
	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}
	l := len(s)
	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}
	for i := startOffset; i < l; i++ {
		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++
		buff.WriteByte(s[i])
	}
	return buff.String()
}
