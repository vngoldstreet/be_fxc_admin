package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func dbMigrations() {
	db_ksc.Migrator().DropTable(&RawMT5Datas{})
	db_ksc.AutoMigrate(&RawMT5Datas{})
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

// func updateWalletForJoinContest(uid uint, cur_balance float64, new_balance float64, change float64, type_id int, contest_id string) (CpsWallets, error) {
// 	newTrans := CpsTransactions{}
// 	newTrans.TypeID = type_id
// 	newTrans.CustomerID = uid
// 	newTrans.CBalance = cur_balance
// 	newTrans.Amount = change
// 	newTrans.NBalance = new_balance
// 	newTrans.StatusID = 1 //Processing
// 	newTrans.ContestID = contest_id
// 	err2 := db_ksc.Create(&newTrans).Error
// 	if err2 != nil {
// 		return CpsWallets{}, err2
// 	}
// 	newWallet := CpsWallets{}
// 	if err := db_ksc.Model(newWallet).Where("customer_id = ?", uid).Update("balance", new_balance).Find(&newWallet).Error; err != nil {
// 		return CpsWallets{}, err
// 	}
// 	if err := db_ksc.Save(&newWallet).Error; err != nil {
// 		return CpsWallets{}, err
// 	}
// 	return newWallet, nil
// }

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

func upLoadFunc(c *gin.Context) {
	fileUpload, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(400, "Bad Request")
		return
	}
	defer fileUpload.Close()

	// Create a directory to store uploaded files
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		c.String(500, "Internal Server Error")
		return
	}
	currentTime := time.Now()
	fileName := fmt.Sprintf("%d%d%d_%dh%d_%s", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), removeSpecialChars(header.Filename))
	// Create a new file on the server
	out, err := os.Create("uploads/" + fileName)
	if err != nil {
		c.String(500, "Internal Server Error")
		return
	}
	defer out.Close()

	// Copy the file content from the form to the file on the server
	_, err = io.Copy(out, fileUpload)
	if err != nil {
		c.String(500, "Internal Server Error")
		return
	}

	// Open and read the uploaded CSV file
	file, err := os.Open("uploads/" + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read and ignore the header row
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	var newUpload []RawMT5Datas
	for i, line := range records {
		if i == 0 {
			continue
		}
		record := strings.Split(line[0], ";")
		if removeSpecialChars(record[5]) != "" {
			balance, _ := strconv.ParseFloat(record[6], 64)
			equity, _ := strconv.ParseFloat(record[7], 64)
			profit, _ := strconv.ParseFloat(record[8], 64)
			floating, _ := strconv.ParseFloat(record[9], 64)
			data := RawMT5Datas{
				Login:      removeSpecialChars(record[0]),
				Name:       removeSpecialChars(record[1]),
				LastName:   removeSpecialChars(record[2]),
				MiddleName: removeSpecialChars(record[3]),
				Email:      removeSpecialChars(record[4]),
				ContestID:  removeSpecialChars(record[5]), //comment
				Balance:    balance,
				Equity:     equity,
				Profit:     profit,
				FloatingPL: floating,
			}
			newUpload = append(newUpload, data)
		}
	}

	listContest := []ListContests{}
	if err := db_ksc.Model(ListContests{}).Select("contest_id").Where("status_id = 1").Find(&listContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := []RawMT5Datas{}

	newDataToCreates := []RawMT5Datas{}
	for _, v := range listContest {
		currentData := []RawMT5Datas{}
		if err := db_ksc.Model(RawMT5Datas{}).Where("contest_id = ?", v.ContestID).Find(&currentData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for _, current := range currentData {
			for _, new := range newUpload {
				fmt.Printf("new.Login: %v - Current: %v\n", new.Login, current.Login)
				if current.Login == new.Login {
					updates := RawMT5Datas{
						Name:       new.Name,
						LastName:   new.LastName,
						MiddleName: new.MiddleName,
						ContestID:  new.ContestID,
						Email:      new.Email,
						Balance:    new.Balance,
						Equity:     new.Equity,
						Profit:     new.Profit,
						FloatingPL: new.FloatingPL,
					}
					if err := db_ksc.Model(&current).Where("id = ?", current.ID).Updates(updates).Error; err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					updated = append(updated, current)
					// idx = true
				}
			}
			// if !idx {
			// 	newCreate := RawMT5Datas{
			// 		Login:      new.Login,
			// 		Name:       new.Name,
			// 		LastName:   new.LastName,
			// 		MiddleName: new.MiddleName,
			// 		ContestID:  new.ContestID,
			// 		Email:      new.Email,
			// 		Balance:    new.Balance,
			// 		Equity:     new.Equity,
			// 		Profit:     new.Profit,
			// 		FloatingPL: new.FloatingPL,
			// 	}
			// 	newDataToCreates = append(newDataToCreates, newCreate)
			// }
		}
	}
	if len(newDataToCreates) > 0 {
		if err := db_ksc.Model(&RawMT5Datas{}).Create(&newDataToCreates).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{
		"contest":     listContest,
		"data_upload": newUpload,
		"update":      updated,
		"new":         newDataToCreates,
	})
}

func updateContestByID(c *gin.Context) {
	var input ListContests
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newContest := ListContests{}
	if err := db_ksc.Where("contest_id = ?", input.ContestID).Find(&newContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newContest.StatusID = input.StatusID
	// newContest.ContestID = GenerateSecureCodeContest(int(newContest.ID))
	if err := db_ksc.Save(&newContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateContest := []Contests{}
	if err := db_ksc.Model(Contests{}).Where("contest_id = ?", input.ContestID).Find(&updateContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, new := range updateContest {
		updates := Contests{
			StatusID: input.StatusID,
		}
		if err := db_ksc.Model(&new).Where("id = ?", new.ID).Updates(updates).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//Delete from redis
	userid := []CpsUsers{}
	db_ksc.Model(CpsUsers{}).Select("id").Find(&userid)
	keysToDelete := []string{}
	for _, v := range userid {
		keysToDelete = append(keysToDelete, setKey(v.ID, db_greetings))
	}
	if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
		fmt.Printf("err Del Redis key: %v\n", err)
	}
	c.JSON(http.StatusOK, gin.H{"data": newContest})
}

func checkTime(time_string string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"

	loc, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Error loading local time zone:", err)
		return time.Now(), err
	}

	parsedTime, err := time.ParseInLocation(layout, time_string, loc)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Now(), err
	}
	return parsedTime, nil
}

func createContest(c *gin.Context) {
	var input CreateContest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime, err := checkTime(input.Start_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endTime, err := checkTime(input.Expired_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newContest := ListContests{}
	newContest.Amount = input.Amount
	newContest.MaximumPerson = input.MaximumPerson
	newContest.Start_at = startTime
	newContest.Expired_at = endTime
	newContest.StartBalance = input.StartBalance
	newContest.EstimatedTime = endTime.Add(24 * time.Hour)
	newContest.StatusID = input.StatusID

	if err := db_ksc.Create(&newContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newContest.ContestID = GenerateSecureCodeContest(int(newContest.ID))
	if err := db_ksc.Save(&newContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	formatMessage := "A new contest has been created: %s\nAmount: %sG\nMaxPerson: %d\nStartAt: %s\nExpiresAt: %s\nStartBalance: %d$\n"
	msg := fmt.Sprintf(formatMessage, newContest.ContestID, NumberToString(int(newContest.Amount), ','), newContest.MaximumPerson, newContest.Start_at.Format("2006-01-02 15:04:05"), newContest.Expired_at.Format("2006-01-02 15:04:05"), newContest.StartBalance)
	if err := SaveToMessages(1, msg); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	//Delete from redis
	userid := []CpsUsers{}
	db_ksc.Model(CpsUsers{}).Select("id").Find(&userid)
	keysToDelete := []string{}
	for _, v := range userid {
		keysToDelete = append(keysToDelete, setKey(v.ID, db_greetings))
	}
	if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
		fmt.Printf("err Del Redis key: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": newContest})
}

func approvalContest(c *gin.Context) {
	var input Contests
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentContest := Contests{}
	if err := db_ksc.Model(&currentContest).Where("customer_id = ? and contest_id = ? and status_id=0", input.CustomerID, input.ContestID).Find(&currentContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentContest.ContestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not found"})
		return
	}
	currentContest.FxID = input.FxID
	currentContest.FxMasterPw = input.FxMasterPw
	currentContest.FxInvesterPw = input.FxInvesterPw
	currentContest.StatusID = 1

	if err := db_ksc.Save(&currentContest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := CpsUsers{}
	if err := db_ksc.Model(&user).Select("name, email").Where("id = ?", input.CustomerID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strLogin := fmt.Sprintf("%d", input.FxID)
	newLeaderBoard := RawMT5Datas{
		Login:     strLogin,
		Name:      user.Name,
		Email:     user.Email,
		ContestID: input.ContestID,
	}

	if err := db_ksc.Save(&newLeaderBoard).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	formatMessage := "Approved a customer's participation in the contest: %s\nCustomer: %s (%d - %s)\nFxID: %d\nFxInvesterPw: %s"
	msg := fmt.Sprintf(formatMessage, input.ContestID, user.Name, input.CustomerID, user.Email, currentContest.FxID, currentContest.FxInvesterPw)

	if err := SaveToMessages(1, msg); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	SendEmailForContest(user.Email, input.ContestID, strLogin, currentContest.FxMasterPw, currentContest.FxInvesterPw)
	//Delete from redis
	keysToDelete := []string{}
	keysToDelete = append(keysToDelete, setKey(input.CustomerID, db_greetings))
	if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
		fmt.Printf("err Del Redis key: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": currentContest})
}

func approvalTransactions(c *gin.Context) {
	var input CpsTransactions
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTrans := CpsTransactions{}

	if err := db_ksc.Model(newTrans).Where("id = ?", input.ID).Find(&newTrans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newWallet := CpsWallets{}
	switch newTrans.TypeID {
	case 1, 3, 5:
		{
			if newTrans.StatusID == 1 {
				wallet := CpsWallets{}
				if err := db_ksc.Where("customer_id = ?", newTrans.CustomerID).Find(&wallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				currentBalance := wallet.Balance
				changeBalance := newTrans.Amount
				newBalance := currentBalance + changeBalance

				if err := db_ksc.Model(newWallet).Where("customer_id = ?", newTrans.CustomerID).Update("balance", newBalance).Find(&newWallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				newTrans.StatusID = 2
				newTrans.CBalance = currentBalance
				newTrans.NBalance = newBalance

				if err := db_ksc.Save(&newTrans).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				user := CpsUsers{}
				if err := db_ksc.Model(&user).Select("id, name, email").Where("id = ?", newTrans.CustomerID).Find(&user).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				type_string := CheckTransactiontype(newTrans.TypeID)
				formatMessage := "Deposit completed: %d\nCustomer: %s (%d - %s)\nType: %s\nBeforeBalance: %sG\nChange: %sG\nNewBalance: %sG\n"
				msg := fmt.Sprintf(formatMessage, newTrans.ID, user.Name, user.ID, user.Email, type_string, NumberToString(int(newTrans.CBalance), ','), NumberToString(int(newTrans.Amount), ','), NumberToString(int(newTrans.NBalance), ','))

				//Delete from redis
				keysToDelete := []string{}
				keysToDelete = append(keysToDelete, setKey(newTrans.CustomerID, db_greetings))
				if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
					fmt.Printf("err Del Redis key: %v\n", err)
				}

				if err := SaveToMessages(2, msg); err != nil {
					fmt.Printf("err: %v\n", err)
				}

				c.JSON(http.StatusOK, gin.H{
					"old_wallet": wallet,
					"new_wallet": newWallet,
				})

			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Transaction does not exist",
				})
			}
		}
	case 2, 4:
		{
			if newTrans.StatusID == 1 {
				newTrans.StatusID = 2
				if err := db_ksc.Save(&newTrans).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				user := CpsUsers{}
				if err := db_ksc.Model(&user).Select("id, name, email").Where("id = ?", newTrans.CustomerID).Find(&user).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				wallet := CpsWallets{}
				if err := db_ksc.Model(wallet).Where("customer_id = ?", user.ID).Find(&wallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				type_string := CheckTransactiontype(newTrans.TypeID)
				formatMessage := "Withdraw complete: %d\nCustomer: %s (%d - %s)\nType: %s\nBalance: %sG"
				msg := fmt.Sprintf(formatMessage, newTrans.ID, user.Name, user.ID, user.Email, type_string, NumberToString(int(wallet.Balance), ','))

				if err := SaveToMessages(2, msg); err != nil {
					fmt.Printf("err: %v\n", err)
				}

				//Delete from redis
				keysToDelete := []string{}
				keysToDelete = append(keysToDelete, setKey(newTrans.CustomerID, db_greetings))
				if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
					fmt.Printf("err Del Redis key: %v\n", err)
				}

				c.JSON(http.StatusOK, gin.H{
					"message": msg,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Transaction does not exist",
				})
			}
		}
	}
}

func cancelTransactions(c *gin.Context) {
	var input CpsTransactions
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTrans := CpsTransactions{}

	if err := db_ksc.Model(newTrans).Where("id = ?", input.ID).Find(&newTrans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newWallet := CpsWallets{}
	switch newTrans.TypeID {
	case 1, 3, 5:
		{
			if newTrans.StatusID == 1 {
				newTrans.StatusID = 3

				errSaveTrans := db_ksc.Save(&newTrans).Error
				if errSaveTrans != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": errSaveTrans.Error()})
					return
				}

				user := CpsUsers{}
				if err := db_ksc.Model(&user).Select("id, name, email").Where("id = ?", newTrans.CustomerID).Find(&user).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				wallet := CpsWallets{}
				if err := db_ksc.Model(wallet).Where("customer_id = ?", user.ID).Find(&wallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				type_string := CheckTransactiontype(newTrans.TypeID)
				formatMessage := "Cancel transaction: %d\nCustomer: %s (%d - %s)\nType: %s\nBalance: %sG"
				msg := fmt.Sprintf(formatMessage, newTrans.ID, user.Name, user.ID, user.Email, type_string, NumberToString(int(wallet.Balance), ','))

				if err := SaveToMessages(2, msg); err != nil {
					fmt.Printf("err: %v\n", err)
				}

				//Delete from redis
				keysToDelete := []string{}
				keysToDelete = append(keysToDelete, setKey(newTrans.CustomerID, db_greetings))
				if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
					fmt.Printf("err Del Redis key: %v\n", err)
				}

				c.JSON(http.StatusOK, gin.H{
					"message": msg,
				})

			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Transaction does not exist",
				})
			}
		}
	case 2, 4:
		{
			if newTrans.StatusID == 1 {
				wallet := CpsWallets{}
				if err := db_ksc.Where("customer_id = ?", newTrans.CustomerID).Find(&wallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				currentBalance := wallet.Balance
				changeBalance := newTrans.Amount
				newBalance := currentBalance + changeBalance

				if err := db_ksc.Model(newWallet).Where("customer_id = ?", newTrans.CustomerID).Update("balance", newBalance).Find(&newWallet).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				newTrans.StatusID = 3
				newTrans.CBalance = currentBalance
				newTrans.NBalance = newBalance

				if err := db_ksc.Save(&newTrans).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				//--
				user := CpsUsers{}
				if err := db_ksc.Model(&user).Select("id, name, email").Where("id = ?", newTrans.CustomerID).Find(&user).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				type_string := CheckTransactiontype(newTrans.TypeID)
				formatMessage := "Cancel transaction: %d\nCustomer: %s (%d - %s)\nType: %s\nBeforeBalance: %sG\nChange: %sG\nNewBalance: %sG\n"
				msg := fmt.Sprintf(formatMessage, newTrans.ID, user.Name, user.ID, user.Email, type_string, NumberToString(int(currentBalance), ','), NumberToString(int(changeBalance), ','), NumberToString(int(newWallet.Balance), ','))

				if err := SaveToMessages(2, msg); err != nil {
					fmt.Printf("err: %v\n", err)
				}

				//Delete from redis
				keysToDelete := []string{}
				keysToDelete = append(keysToDelete, setKey(newTrans.CustomerID, db_greetings))
				if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
					fmt.Printf("err Del Redis key: %v\n", err)
				}

				//--
				c.JSON(http.StatusOK, gin.H{
					"old_wallet": wallet,
					"new_wallet": newWallet,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Transaction does not exist",
				})
			}
		}
	}
}

func createTransactions(c *gin.Context) {
	var input CpsTransactions
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//create new transaction
	newTrans := CpsTransactions{}
	newTrans.TypeID = input.TypeID
	newTrans.CustomerID = input.CustomerID
	newTrans.Amount = input.Amount
	newTrans.StatusID = 1 //Processing

	if err := db_ksc.Create(&newTrans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resTransaction := []CpsTransactions{}
	if err := db_ksc.Model(CpsTransactions{}).Where("customer_id = ?", input.CustomerID).Order("id desc").Find(&resTransaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get wallet to send message to telegram
	resWallet := CpsWallets{}
	if err := db_ksc.Model(CpsWallets{}).Where("customer_id = ?", input.CustomerID).Find(&resWallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//get users
	resUser := CpsUsers{}
	if err := db_ksc.Model(CpsUsers{}).Where("id = ?", input.CustomerID).Find(&resUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type_string := CheckTransactiontype(newTrans.TypeID)
	formatMessage := "Create a transaction by admin: %d\nCustomer: %s (%d - %s)\nType: %s\nAmount: %s\nBalance: %sG"
	msg := fmt.Sprintf(formatMessage, newTrans.ID, resUser.Name, resUser.ID, resUser.Email, type_string, NumberToString(int(newTrans.Amount), ','), NumberToString(int(resWallet.Balance), ','))

	if err := SaveToMessages(2, msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Delete from redis
	keysToDelete := []string{}
	keysToDelete = append(keysToDelete, setKey(input.CustomerID, db_greetings))
	if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
		fmt.Printf("err Del Redis key: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func getTransactions(c *gin.Context) {
	trans := []CpsAdminTransactions{}
	selectPromp := `cps_transactions.id as id,
                  cps_transactions.customer_id as customer_id,
                  cps_transactions.amount as amount,
                  cps_transactions.status_id as status_id,
				  cps_transactions.type_id as type_id,
                  cps_users.name as name,
                  cps_users.phone as phone,
                  cps_users.email as email,
                  cps_transactions.created_at as created_at,
                  cps_transactions.updated_at as updated_at`
	if err := db_ksc.Model(&CpsTransactions{}).Select(selectPromp).Joins("INNER JOIN cps_users on cps_transactions.customer_id = cps_users.id").Where("status_id = 1 and type_id <> 4").Order("cps_transactions.id desc").Find(&trans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trans})
}

func getCompetitionRequest(c *gin.Context) {
	trans := []CpsAdminTransactions{}
	selectPromp := `cps_transactions.id as id,
                  cps_transactions.customer_id as customer_id,
                  cps_transactions.amount as amount,
                  cps_transactions.status_id as status_id,
				  cps_transactions.type_id as type_id,
				  cps_transactions.contest_id as contest_id,
                  cps_users.name as name,
				  cps_users.id as customer_id,
                  cps_users.phone as phone,
                  cps_users.email as email,
                  cps_transactions.created_at as created_at,
                  cps_transactions.updated_at as updated_at`
	if err := db_ksc.Model(&CpsTransactions{}).Select(selectPromp).Joins("INNER JOIN cps_users on cps_transactions.customer_id = cps_users.id").Where("status_id = 1 and type_id = 4").Order("cps_transactions.id desc").Find(&trans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trans})
}

func getCompetitionRequestHistory(c *gin.Context) {
	trans := []CpsAdminTransactions{}
	selectPromp := `cps_transactions.id as id,
                  cps_transactions.customer_id as customer_id,
                  cps_transactions.amount as amount,
                  cps_transactions.status_id as status_id,
				  cps_transactions.type_id as type_id,
				  cps_transactions.contest_id as contest_id,
                  cps_users.name as name,
				  cps_users.id as customer_id,
                  cps_users.phone as phone,
                  cps_users.email as email,
                  cps_transactions.created_at as created_at,
                  cps_transactions.updated_at as updated_at`
	if err := db_ksc.Model(&CpsTransactions{}).Select(selectPromp).Joins("INNER JOIN cps_users on cps_transactions.customer_id = cps_users.id").Where("status_id <> 1 and type_id = 4").Order("cps_transactions.id desc").Find(&trans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trans})
}

func getHistoryTransactions(c *gin.Context) {
	trans := []CpsAdminTransactions{}
	// if err := db_ksc.Model(CpsTransactions{}).Where("status_id = 1 and type_id <> 4").Find(&trans).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	selectPromp := `cps_transactions.id as id,
                  cps_transactions.customer_id as customer_id,
                  cps_transactions.amount as amount,
                  cps_transactions.status_id as status_id,
				  cps_transactions.type_id as type_id,
                  cps_users.name as name,
                  cps_users.phone as phone,
                  cps_users.email as email,
                  cps_transactions.created_at as created_at,
                  cps_transactions.updated_at as updated_at`
	if err := db_ksc.Model(&CpsTransactions{}).Select(selectPromp).Joins("INNER JOIN cps_users on cps_transactions.customer_id = cps_users.id").Where("status_id <> 1 and type_id <> 4").Order("cps_transactions.id desc").Find(&trans).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trans})
}

func getContestList(c *gin.Context) {
	datas := []ListContests{}
	if err := db_ksc.Model(ListContests{}).Where("status_id in (0,1)").Order("id desc").Find(&datas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": datas})
}

func getHistoryContestList(c *gin.Context) {
	datas := []ListContests{}
	if err := db_ksc.Model(ListContests{}).Where("status_id > 1").Order("id desc").Find(&datas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": datas})
}

func getReviewLists(c *gin.Context) {
	datas := []CpsAdminReviews{}
	selectPromp := `cps_reviews.id as id,
                  cps_reviews.customer_id as customer_id,
                  cps_reviews.image_front as image_front,
                  cps_reviews.image_back as image_back,
				  cps_reviews.status as status,
                  cps_users.name as name,
                  cps_users.phone as phone,
                  cps_users.email as email,
                  cps_reviews.created_at as created_at,
                  cps_reviews.updated_at as updated_at`
	if err := db_ksc.Model(&CpsReviews{}).Select(selectPromp).Joins("INNER JOIN cps_users on cps_reviews.customer_id = cps_users.id").Where("status <> ?", "done").Order("cps_reviews.id desc").Find(&datas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": datas})
}

func updateReviewLists(c *gin.Context) {
	var input CpsReviews
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datas := CpsReviews{}
	if err := db_ksc.Model(CpsReviews{}).Where("customer_id = ?", input.CustomerID).Order("id desc").Find(&datas).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := CpsUsers{}
	if input.Status == "done" {
		datas.Status = "done"
		if err := db_ksc.Save(&datas).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db_ksc.Model(&user).Where("id = ?", input.CustomerID).Find(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.InReview = "done"
		if err := db_ksc.Save(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": datas, "user": user})
}
