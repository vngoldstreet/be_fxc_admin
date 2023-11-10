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
	r.GET("/transactions-history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transactions-history.html", gin.H{
			"title": "History of Transaction",
		})
	})
	r.GET("/competitions", func(c *gin.Context) {
		c.HTML(http.StatusOK, "competitions.html", gin.H{
			"title": "Competitions",
		})
	})
	r.GET("/competitions-history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "competitions-history.html", gin.H{
			"title": "History of Competitions",
		})
	})
	r.GET("/contests", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contests.html", gin.H{
			"title": "Contests",
		})
	})
	r.GET("/contests-history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contests-history.html", gin.H{
			"title": "History of Contests",
		})
	})
	r.GET("/customer-review", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inreviews.html", gin.H{
			"title": "Customer inreviews",
		})
	})
	r.GET("/uploader", func(c *gin.Context) {
		c.HTML(http.StatusOK, "leaderboard.html", gin.H{
			"title": "Upload",
		})
	})
	r.GET("/create-post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create-post.html", gin.H{
			"title": "Create Post",
		})
	})
	r.GET("/update-post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update-posts.html", gin.H{
			"title": "Update Post",
		})
	})
	public := r.Group("/public")
	public.POST("/register", Register)
	public.POST("/login", Login)

	public.GET("/posts", getPosts)
	public.GET("/image", getImage)
	public.GET("/post-by-url", getPostByUrl)

	private := r.Group("/auth")
	private.Use(JwtAuthMiddleware())
	private.POST("/upload-csv", upLoadFunc)
	private.POST("/upload-old-leaderboard", upLoadOldLeaderboard)
	private.POST("/create-contest", createContest)
	private.POST("/create-transaction", createTransactions)
	private.POST("/update-contest-id", updateContestByID)                              //done
	private.POST("/contest-approval", approvalContest)                                 //done
	private.POST("/rejoin-contest-approval", approvalRejoinContest)                    //done
	private.POST("/admin-transaction", approvalTransactions)                           //done
	private.POST("/cancel-transaction", cancelTransactions)                            //done
	private.GET("/get-transaction-list", getTransactions)                              //done
	private.GET("/get-history-transaction-list", getHistoryTransactions)               //done
	private.GET("/get-competition-request-list", getCompetitionRequest)                //done
	private.GET("/get-history-competition-request-list", getCompetitionRequestHistory) //done
	private.GET("/get-contest-list", getContestList)                                   //done
	private.GET("/get-history-contest-list", getHistoryContestList)                    //done
	private.GET("/get-review-list", getReviewLists)                                    //done
	private.POST("/update-review-list", updateReviewLists)                             //done
	private.POST("/create-post", postDatas)                                            //done
	private.POST("/update-post", updatePost)                                           //done

	r.Run(":8081")
}
