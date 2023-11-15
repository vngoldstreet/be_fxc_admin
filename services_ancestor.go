package main

import (
	"fmt"
	"net/http"

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
