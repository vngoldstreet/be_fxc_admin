package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/src/assets", "./src/assets")
	r.LoadHTMLGlob("src/html/*")
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "authentication-login.html", gin.H{
			"title": "Login",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transactions.html", gin.H{
			"title": "Transactions",
		})
	})
	r.GET("/transaction-history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transaction-history.html", gin.H{
			"title": "History",
		})
	})
	public := r.Group("/public")
	// public.POST("/register", Register)
	public.POST("/login", Login)

	private := r.Group("/auth")
	private.Use(JwtAuthMiddleware())
	private.POST("/upload-csv", upLoadFunc)
	private.POST("/update-contest-id", updateContestByID)
	private.POST("/create-contest", createContest)
	private.POST("/contest-approval", approvalContest)
	private.POST("/admin-transaction", approvalTransactions) //done
	private.POST("/cancel-transaction", cancelTransactions)  //done
	private.POST("/create-transaction", createTransactions)
	private.GET("/get-transaction-list", getTransactions) //done
	private.GET("/get-history-transaction-list", getHistoryTransactions)
	private.GET("/get-contest-list", getContestList)
	private.GET("/get-history-contest-list", getHistoryContestList)
	private.GET("/get-review-list", getReviewLists)
	private.POST("/update-review-list", updateReviewLists)
	r.Run(":8081")

}
