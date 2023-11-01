package main

import (
	"regexp"
	"time"

	"gorm.io/gorm"
)

type CpsUsers struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	PartnerCode string `json:"code"`
	Image       string `json:"image"`
	Description string `json:"description"`
	RefLink     string `json:"ref_link"`
	InReview    string `json:"inreview" gorm:"default:not_yet"`
}

type CpsReviews struct {
	gorm.Model
	CustomerID uint   `json:"customer_id"`
	ImageFront string `json:"image_front"`
	ImageBack  string `json:"image_back"`
	Status     string `json:"status"`
}

type CpsAdminReviews struct {
	gorm.Model
	CustomerID uint   `json:"customer_id"`
	ImageFront string `json:"image_front"`
	ImageBack  string `json:"image_back"`
	Status     string `json:"status"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

// Users
type CpsAdmins struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CpsWallets struct {
	gorm.Model
	CustomerID uint    `json:"customer_id"`
	Balance    float64 `json:"balance"`
	LastChange float64 `json:"last_change"`
}

// Partners
type CpsPartners struct {
	gorm.Model
	ParentID int    `json:"parent_id"`
	ChildID  int    `json:"child_id"`
	Path     string `json:"path"`
}

// Transaction
type CpsTransactions struct {
	gorm.Model
	TypeID        int     `json:"type_id"`
	CustomerID    uint    `json:"customer_id"`
	CBalance      float64 `json:"c_balance"`
	Amount        float64 `json:"amount"`
	NBalance      float64 `json:"n_balance"`
	PaymentMethob int     `json:"payment_methob"`
	PaymentGate   int     `json:"payment_gate"`
	StatusID      int     `json:"status_id"`
	ParentID      int     `json:"parent_id"`
	ContestID     string  `json:"contest_id"`
}

type CpsAdminTransactions struct {
	gorm.Model
	TypeID        int     `json:"type_id"`
	CustomerID    uint    `json:"customer_id"`
	CBalance      float64 `json:"c_balance"`
	Amount        float64 `json:"amount"`
	NBalance      float64 `json:"n_balance"`
	PaymentMethob int     `json:"payment_methob"`
	PaymentGate   int     `json:"payment_gate"`
	StatusID      int     `json:"status_id"`
	ParentID      int     `json:"parent_id"`
	ContestID     string  `json:"contest_id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	HolderName    string  `json:"holder_name"`
	HolderNumber  string  `json:"holder_number"`
	BankName      string  `json:"bank_name"`
}

type CpsTransactionTypes struct {
	gorm.Model
	Type        int    `json:"type"`
	Description string `json:"description"`
}

type CpsTransactionStatus struct {
	gorm.Model
	Status      int    `json:"status"`
	Description string `json:"description"`
}

type CpsPaymentMethobs struct {
	gorm.Model
	CustomerID   uint   `json:"customer_id"`
	HolderName   string `json:"holder_name"`
	HolderNumber string `json:"holder_number"`
	BankName     string `json:"bank_name"`
	IsCard       int    `json:"is_card" gorm:"default:0"`
}

type CpsPaymentGates struct {
	gorm.Model
	Status      int    `json:"status"`
	Description string `json:"description"`
}

type CpsNotifications struct {
	gorm.Model
	CustomerID uint   `json:"customer_id"`
	Type       int    `json:"type"`
	Message    string `json:"message"`
	IsSent     int    `json:"is_send" gorm:"default:0"`
}

// Contest
type ListContests struct {
	gorm.Model
	ContestID     string    `json:"contest_id"`
	Amount        float64   `json:"amount"`
	MaximumPerson int       `json:"maximum_person"`
	CurrentPerson int       `json:"current_person" gorm:"default:0"`
	Start_at      time.Time `json:"start_at"`
	Expired_at    time.Time `json:"expired_at"`
	StartBalance  int       `json:"start_balance"`
	EstimatedTime time.Time `json:"estimate_time"`
	StatusID      int       `json:"status_id" gorm:"default:0"`
}

type CreateContest struct {
	ContestID     string  `json:"contest_id"`
	Amount        float64 `json:"amount"`
	MaximumPerson int     `json:"maximum_person"`
	CurrentPerson int     `json:"current_person" gorm:"default:0"`
	Start_at      string  `json:"start_at"`
	Expired_at    string  `json:"expired_at"`
	StartBalance  int     `json:"start_balance"`
	EstimatedTime string  `json:"estimate_time"`
	StatusID      int     `json:"status_id" gorm:"default:0"`
}

type Contests struct {
	gorm.Model
	ContestID    string `json:"contest_id"`
	CustomerID   uint   `json:"customer_id"`
	FxID         uint   `json:"fx_id"`
	FxMasterPw   string `json:"fx_master_pw"`
	FxInvesterPw string `json:"fx_invester_pw"`
	StatusID     int    `json:"status_id" gorm:"default:0"`
}

type ContestInfos struct {
	ContestID     string    `json:"contest_id"`
	Amount        float64   `json:"amount"`
	MaximumPerson int       `json:"maximum_person"`
	CurrentPerson int       `json:"current_person"`
	Start_at      time.Time `json:"start_at"`
	Expired_at    time.Time `json:"expired_at"`
	StartBalance  int       `json:"start_balance"`
	StatusID      int       `json:"status_id" gorm:"default:0"`
	CustomerID    uint      `json:"customer_id"`
	FxID          uint      `json:"fx_id"`
	FxMasterPw    string    `json:"fx_master_pw"`
	FxInvesterPw  string    `json:"fx_invester_pw"`
}

type LeaderBoards struct {
	gorm.Model
	ContestID      string  `json:"contest_id"`
	CustomerID     uint    `json:"customer_id"`
	StartBalance   float64 `json:"start_balance"`
	CurrentBalance float64 `json:"current_balance"`
	CurrentEquity  float64 `json:"current_equity"`
	FloatingPL     float64 `json:"pnl"`
}

type CpsMessages struct {
	gorm.Model
	TypeID  int    `json:"type_id"`
	Message string `json:"message"`
	IsSent  int    `json:"is_sent" gorm:"default:0"`
}

func removeSpecialChars(input string) string {
	allowedChars := regexp.MustCompile(`[a-zA-Z0-9@,.\sáàảãạăắằẳẵặâấầẩẫậéèẻẽẹêếềểễệíìỉĩịóòỏõọôốồổỗộơớờởỡợúùủũụưứừửữựýỳỷỹỵ\-+/*]+`)

	// Loại bỏ tất cả các ký tự khác
	replaced := allowedChars.FindAllString(input, -1)
	result := ""
	for _, s := range replaced {
		result += s
	}

	return result
}

// func removeSpecialChars(input string) string {
// 	regex := regexp.MustCompile("[^a-zA-Z0-9@.,]")

// 	// Sử dụng ReplaceAllString để thay thế tất cả các ký tự đặc biệt bằng dấu trống
// 	result := regex.ReplaceAllString(input, "")

// 	return result
// }

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
