package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RawMT5Datas struct {
	gorm.Model
	Login      string  `json:"login"`
	Name       string  `json:"name"`
	LastName   string  `json:"last_name"`
	MiddleName string  `json:"middle_name"`
	ContestID  string  `json:"contest_id"`
	Email      string  `json:"email"`
	Balance    float64 `json:"balance"`
	Equity     float64 `json:"equity"`
	Profit     float64 `json:"profit"`
	FloatingPL float64 `json:"floating"`
}

func main() {
	// dbMigrations()
	// db_ksc.Migrator().DropTable(RawMT5Datas{})
	// db_ksc.AutoMigrate(RawMT5Datas{})
	r := gin.Default()
	r.Static("/src/assets", "./src/assets")
	r.LoadHTMLGlob("html/*")
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title": "Upload",
		})
	})

	public := r.Group("/public")
	public.POST("/register", Register)
	public.POST("/login", Login)
	public.GET("/send_message", func(c *gin.Context) {
		if err := GetAndSendMessageFromDb(); err != nil {
			c.JSON(http.StatusOK, gin.H{"mess": err})
		}
		c.JSON(http.StatusOK, gin.H{"mess": "success"})
	})
	private := r.Group("/auth")
	// private.Use(JwtAuthMiddleware())
	r.POST("/upload-csv", func(c *gin.Context) {
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

		var datas []RawMT5Datas
		for i, line := range records {
			if i < 2 {
				continue
			}
			if i == len(records)-1 {
				continue
			}

			record := strings.Split(line[0], ";")
			if removeSpecialChars(record[4]) != "" {
				balance, _ := strconv.ParseFloat(removeSpecialChars(record[6]), 64)
				equity, _ := strconv.ParseFloat(removeSpecialChars(record[7]), 64)
				profit, _ := strconv.ParseFloat(removeSpecialChars(record[8]), 64)
				floating, _ := strconv.ParseFloat(removeSpecialChars(record[9]), 64)
				data := RawMT5Datas{
					Login:      removeSpecialChars(record[0]),
					Name:       removeSpecialChars(record[1]),
					LastName:   removeSpecialChars(record[2]),
					MiddleName: removeSpecialChars(record[3]),
					ContestID:  removeSpecialChars(record[4]),
					Email:      removeSpecialChars(record[5]),
					Balance:    balance,
					Equity:     equity,
					Profit:     profit,
					FloatingPL: floating,
				}
				datas = append(datas, data)
			}
		}
		db_ksc.Save(&datas)

		//Delete from redis
		userid := []CpsUsers{}
		db_ksc.Model(CpsUsers{}).Select("id").Find(&userid)
		keysToDelete := []string{}
		for _, v := range userid {
			keysToDelete = append(keysToDelete, setKey(v.ID, db_leaderboard))
		}
		if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
			fmt.Printf("err Del Redis key: %v\n", err)
		}

		c.JSON(200, datas)
	})
	private.POST("/update-contest-id", func(c *gin.Context) {
		var input ListContests
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newContest := ListContests{}
		if err := db_ksc.Where("id = ?", input.ID).Find(&newContest).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newContest.ContestID = GenerateSecureCodeContest(int(newContest.ID))
		if err := db_ksc.Save(&newContest).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
	})
	private.POST("/create-contest", func(c *gin.Context) {
		var input ListContests
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeCurrent := time.Now()

		newContest := ListContests{}
		newContest.Amount = input.Amount
		newContest.MaximumPerson = input.MaximumPerson
		newContest.Start_at = timeCurrent
		newContest.Expired_at = timeCurrent
		newContest.StartBalance = input.StartBalance
		newContest.EstimatedTime = timeCurrent
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
	})
	private.POST("/contest-approval", func(c *gin.Context) {
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
		formatMessage := "Approved a customer's participation in the contest: %s\nCustomer: %s (%d - %s)\nFxID: %d\nFxInvesterPw: %s"
		msg := fmt.Sprintf(formatMessage, input.ContestID, user.Name, input.CustomerID, user.Email, currentContest.FxID, currentContest.FxInvesterPw)

		if err := SaveToMessages(1, msg); err != nil {
			fmt.Printf("err: %v\n", err)
		}

		//Delete from redis
		keysToDelete := []string{}
		keysToDelete = append(keysToDelete, setKey(input.CustomerID, db_greetings))
		if _, err := rdb.Del(context.Background(), keysToDelete...).Result(); err != nil {
			fmt.Printf("err Del Redis key: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"data": currentContest})
	})
	private.POST("/admin-transaction", func(c *gin.Context) {
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
	})
	private.POST("/cancel-transaction", func(c *gin.Context) {
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
	})
	r.Run(":8081")

}
