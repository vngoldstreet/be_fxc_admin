package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db_ksc *gorm.DB
var rdb *redis.Client

func init() {
	connectDatabase()

	// SendEmailForContest("vietvufx@gmail.com", "abchdwr", "8008000", "khongshochay", "khongsochay")
	// SendEmailForRegister("vietvufx@gmail.com", "8008000", "khongshochay")
	dbMigrations()
	setupLogger()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Thay thế bằng địa chỉ Redis thực tế
		Password: "",               // Mật khẩu (nếu có)
		DB:       0,                // Chọn cơ sở dữ liệu
	})
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Ping err: %v\n", err)
	}
	fmt.Printf("ping: %v\n", ping)
}

func connectDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot connect to database ", DbHost)
		log.Fatal("connection error:", err)
	} else {
		fmt.Println("Connected: ", DbHost)
	}
	db_ksc = db
}

func setKey(id uint, dbstring string) string {
	return fmt.Sprintf("%s_id:%d", dbstring, id)
}

var db_greetings string = "db:greetings"

// var db_leaderboard string = "db:leaderboards"

var logrusApp = logrus.New()

func setupLogger() {
	logrusApp.SetFormatter(&logrus.JSONFormatter{}) // Định dạng log JSON (tùy chọn)

	// Mở tệp log để lưu
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrusApp.SetOutput(file)
	} else {
		logrusApp.Info("Failed to log to file, using default stderr")
	}

	// Cấu hình logger cho Gin
	gin.SetMode(gin.ReleaseMode) // Đặt chế độ Gin (ReleaseMode hoặc DebugMode)
	gin.DefaultWriter = logrusApp.Writer()
}

// func Logger(c *gin.Context) {
// 	log := logrus.New()

// 	// Logging thông tin về yêu cầu
// 	log.WithFields(logrus.Fields{
// 		"clientIP": c.ClientIP(),
// 		"method":   c.Request.Method,
// 		"path":     c.FullPath(),
// 	}).Info("Request")

// 	// Thực hiện xử lý

// 	// Logging thông tin về phản hồi
// 	log.WithFields(logrus.Fields{
// 		"clientIP": c.ClientIP(),
// 		"status":   c.Writer.Status(),
// 	}).Info("Response")
// }
