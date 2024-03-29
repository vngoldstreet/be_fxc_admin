package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	requestTracker sync.Map
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			fullPath := c.FullPath()

			// Tạo key duy nhất cho mỗi full path và IP
			key := fmt.Sprintf("%s-%s", fullPath, c.ClientIP())

			// Kiểm tra xem request có được xử lý trong 1s không
			if _, exists := requestTracker.LoadOrStore(key, true); exists {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests for the same full path"})
				return
			}

			// Thiết lập thời gian hết hạn sau 1s để giải phóng bộ nhớ
			time.AfterFunc(3*time.Second, func() {
				requestTracker.Delete(key)
			})
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	// clientTracker := NewClientTracker()
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
	r.GET("/partners", func(c *gin.Context) {
		// commissions := []CommissionLevels{}
		cookie, err := c.Cookie("token")
		if err != nil {
			c.HTML(http.StatusOK, "partner.html", gin.H{
				"title": "Partner",
			})
			return
		}
		uuid, err := ExtractTokenIDWithString(cookie)
		if err != nil || uuid != 1 {
			c.HTML(http.StatusOK, "partner.html", gin.H{
				"title": "Partner",
			})
			return
		}

		commissions := []ResponseCommissionLevels{}
		selectPromp := `commission_levels.type_id as type_id,
					commission_levels.partner_id as partner_id,
					commission_levels.level_1 as level_1,
					commission_levels.level_2 as level_2,
					commission_levels.level_3 as level_3,
					commission_levels.level_4 as level_4,
					commission_levels.level_5 as level_5,
					commission_levels.commission_1 as commission_1,
					commission_levels.commission_2 as commission_2,
					commission_levels.commission_3 as commission_3,
					commission_levels.commission_4 as commission_4,
					commission_levels.commission_5 as commission_5,
					cps_users.name as name,
					cps_users.phone as phone,
					cps_users.email as email
				  `
		if err := db_ksc.Model(&CommissionLevels{}).Select(selectPromp).Joins("INNER JOIN cps_users on commission_levels.partner_id = cps_users.id").Find(&commissions).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			// conn.WriteMessage(msgType, []byte("dataContestLists err: "+err.Error()))
			return
		}

		c.HTML(http.StatusOK, "partner.html", gin.H{
			"title":      "Partner",
			"commission": commissions,
		})
	})
	r.GET("/create-post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create-post.html", gin.H{
			"title": "Create Post",
		})
	})
	r.GET("/update-post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update-post.html", gin.H{
			"title": "Update Post",
		})
	})

	public := r.Group("/public")
	public.POST("/register", Register)
	public.POST("/login", Login)
	public.POST("/reset-password", resetPassword) //done
	public.GET("/posts", getPosts)
	public.GET("/all-of-posts", getAllPosts)
	public.GET("/image/:url", getImage)
	public.GET("/post-by-url", getPostByUrl)
	public.GET("/post-by-id", getPostByID)

	private := r.Group("/auth")
	private.Use(JwtAuthMiddleware())
	private.POST("/upload-csv", upLoadFunc)
	private.POST("/upload-old-leaderboard", upLoadOldLeaderboard)
	private.POST("/create-contest", createContest)
	private.POST("/create-transaction", createTransactions)
	private.POST("/update-contest-id", updateContestByID) //done

	private.Use(RateLimitMiddleware())
	private.POST("/contest-approval", approvalContest) //done
	private.GET("/contest/get-account-store", contestGetAccountStore)
	private.GET("/contest/get-current-contest", contestGetCurrentContest)
	private.POST("/contest/send-email", sendEmailFromAdmin)

	private.POST("/rejoin-contest-approval", approvalRejoinContest) //done

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
	private.POST("/delete-post", deletePost)                                           //done
	private.POST("/active-partner", activePartner)                                     //done

	private.POST("/create-store", CreateStore) //done
	r.Run(":8081")
}
