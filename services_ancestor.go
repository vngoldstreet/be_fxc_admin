package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Ancestors struct {
	gorm.Model
	ParentID     uint   `json:"parent_id"`
	PartnerID    uint   `json:"partner_id"`
	AncestorPath string `json:"ancestor_path"`
}

type Partners struct {
	ID         int     `json:"id"`
	CustomerID uint    `json:"customer_id"`
	IsPartner  int     `json:"is_partner"`
	Commission float64 `json:"commission"`
}

func activePartner(c *gin.Context) {
	var input Ancestors
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db_ksc.Begin()
	newAncestor := Ancestors{
		ParentID:     input.ParentID,
		PartnerID:    input.PartnerID,
		AncestorPath: fmt.Sprintf("%d_%d", input.ParentID, input.PartnerID),
	}
	if err := tx.Create(&newAncestor).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activePartner := Partners{
		CustomerID: input.PartnerID,
		IsPartner:  1,
		Commission: 0.3,
	}

	if err := tx.Create(&activePartner).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status":   "Success",
		"ancestor": newAncestor,
		"partner":  activePartner,
	})
}

// 174
func CalculateCommission(transaction CpsTransactions) error {
	db_ksc.AutoMigrate(&Commissions{})
	tx := db_ksc.Begin()
	//Check phả hệ của partner
	ancestor := Ancestors{}
	if err := tx.Model(&Ancestors{}).Where("partner_id = ?", transaction.CustomerID).First(&ancestor).Error; err != nil {
		tx.Rollback()
		return err
	}

	arrString := strings.Split(ancestor.AncestorPath, "_")
	arrInt := []int{}
	for _, v := range arrString {
		num, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return err
		}
		arrInt = append(arrInt, num)
	}
	index := findIndex(arrInt, int(transaction.CustomerID))

	for i, v := range arrInt {
		if i > index {
			continue
		}
		newCommission := Commissions{
			TransactionID:   int(transaction.ID),
			TransactionType: transaction.TypeID,
			ParentID:        v,
			ContestID:       transaction.ContestID,
		}
		if err := tx.Model(&Commissions{}).Create(&newCommission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func findIndex(arr []int, target int) int {
	for i, value := range arr {
		if value == target {
			return i // Trả về vị trí của số nguyên trong mảng
		}
	}
	return -1 // Trả về -1 nếu số nguyên không được tìm thấy trong mảng
}

type Commissions struct {
	gorm.Model
	TransactionID   int    `json:"transaction_id"`
	TransactionType int    `json:"transaction_type"`
	ParentID        int    `json:"parent_id"`
	ContestID       string `json:"contest_id"`
}
