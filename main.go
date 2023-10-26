package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// r.Static("/src/assets", "./src/assets")
	// r.LoadHTMLGlob("html/*")
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
	private.Use(JwtAuthMiddleware())
	private.POST("/upload-csv", upLoadFunc)
	private.POST("/update-contest-id", updateContestByID)
	private.POST("/create-contest", createContest)
	private.POST("/contest-approval", approvalContest)
	private.POST("/admin-transaction", approvalTransactions)
	private.POST("/cancel-transaction", cancelTransactions)
	private.POST("/create-transaction", createTransactions)
	private.GET("/get-transaction-list", getTransactions)
	private.GET("/get-history-transaction-list", getHistoryTransactions)
	private.GET("/get-contest-list", getContestList)
	private.GET("/get-history-contest-list", getHistoryContestList)
	private.GET("/get-review-list", getReviewLists)
	private.POST("/update-review-list", updateReviewLists)
	r.Run(":8081")

}
