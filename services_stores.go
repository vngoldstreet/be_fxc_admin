package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountStores struct {
	gorm.Model
	FxID         string `json:"fx_id"`
	FxMasterPw   string `json:"fx_master_pw"`
	FxInvesterPw string `json:"fx_invester_pw"`
	TypeID       int    `json:"type_id"`                    // 1=silver;2=gold;3=platinum
	StatusID     int    `json:"status_id" gorm:"default:0"` //0 empty;1 active;
}

func CreateStore(c *gin.Context) {
	var input AccountStores
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db_ksc.Begin()
	currentData := []AccountStores{}
	if err := tx.Where("fx_id = ?", input.FxID).Find(&currentData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(currentData) > 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Dữ liệu đã tồn tại",
		})
		return
	}
	newStore := AccountStores{
		FxID:         input.FxID,
		FxMasterPw:   input.FxMasterPw,
		FxInvesterPw: input.FxInvesterPw,
		TypeID:       input.TypeID,
		StatusID:     0,
	}
	if err := tx.Create(&newStore).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newStore,
	})
}
