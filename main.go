package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// type RequestInfo struct {
// 	IP  string
// 	UID uint
// }

// type ClientTracker struct {
// 	Counters   map[RequestInfo]int
// 	Mutex      sync.Mutex
// 	LockTimers map[RequestInfo]*time.Timer
// }

// func NewClientTracker() *ClientTracker {
// 	return &ClientTracker{
// 		Counters:   make(map[RequestInfo]int),
// 		Mutex:      sync.Mutex{},
// 		LockTimers: make(map[RequestInfo]*time.Timer),
// 	}
// }

// func TokenAndIPFilterMiddleware(tracker *ClientTracker) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if c.Request.Method == http.MethodPost {
// 			ip := c.ClientIP()
// 			uuid, err := ExtractTokenID(c)
// 			if err != nil {
// 				return
// 			}

// 			requestInfo := RequestInfo{
// 				IP:  ip,
// 				UID: uuid,
// 			}

// 			tracker.Mutex.Lock()
// 			defer tracker.Mutex.Unlock()

// 			tracker.Counters[requestInfo]++
// 			if tracker.Counters[requestInfo] > 10 {
// 				if _, ok := tracker.LockTimers[requestInfo]; !ok {
// 					fmt.Println("Client locked. UID:", requestInfo.UID)

// 					// tx := db_ksc.Begin()
// 					// resUser := CpsUsers{}
// 					// if err := tx.Where("id = ?", uuid).First(&resUser).Error; err != nil {
// 					// 	tx.Rollback()
// 					// 	return
// 					// }

// 					// msg := fmt.Sprintf("Too many requests from this client: %d - %s (%s). Locked!", resUser.ID, resUser.Name, resUser.Email)
// 					// SendMessageToAccountGroup(msg)

// 					// tx.Commit()

// 					tracker.LockTimers[requestInfo] = time.AfterFunc(3*time.Second, func() {
// 						fmt.Println("Client unlocked. UID:", requestInfo.UID)
// 						delete(tracker.LockTimers, requestInfo)
// 						tracker.Counters[requestInfo] = 0
// 					})
// 				}

// 				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests from this client"})
// 				return
// 			}
// 		}
// 		c.Next()
// 	}
// }

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
	// private.Use(TokenAndIPFilterMiddleware(clientTracker))
	private.POST("/upload-csv", upLoadFunc)
	private.POST("/upload-old-leaderboard", upLoadOldLeaderboard)
	private.POST("/create-contest", createContest)
	private.POST("/create-transaction", createTransactions)
	private.POST("/update-contest-id", updateContestByID) //done

	private.POST("/contest-approval", approvalContest) //done
	private.GET("/contest/get-account-store", contestGetAccountStore)
	private.GET("/contest/get-current-contest", contestGetCurrentContest)
	private.POST("/contest/send-email", sendEmailFromAdmin)

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
	private.POST("/delete-post", deletePost)                                           //done
	private.POST("/active-partner", activePartner)                                     //done

	private.POST("/create-store", CreateStore) //done
	r.Run(":8081")
}
