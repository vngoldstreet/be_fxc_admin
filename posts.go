package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getPosts(c *gin.Context) {
	type_post := c.Query("type")
	length, _ := strconv.Atoi(c.Query("length"))
	datas := []Posts{}
	resp := db_ksc.Where("type = ?", type_post).Order("id desc").Limit(length).Find(&datas)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"count":  resp.RowsAffected,
		"data":   datas,
	})
}

func getAllPosts(c *gin.Context) {
	datas := []Posts{}
	resp := db_ksc.Select("title,description,url,created_at,updated_at").Order("id desc").Find(&datas)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"count":  resp.RowsAffected,
		"data":   datas,
	})
}

func getImage(c *gin.Context) {
	url := c.Param("url")
	datas := Posts{}
	db_ksc.Select("thumb").Where("url=?", url).Find(&datas)
	i := strings.Index(datas.Thumb, ",")
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(datas.Thumb[i+1:]))
	c.Header("Content-Type", "image/png, image/jpg")
	io.Copy(c.Writer, dec)
}

func postDatas(c *gin.Context) {
	var input Posts
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := Posts{
		Title:       input.Title,
		Description: input.Description,
		Content:     input.Content,
		Thumb:       input.Thumb,
		Type:        input.Type,
		Tag:         input.Tag,
		Viewer:      input.Viewer,
		Url:         input.Url,
		Keyword:     input.Keyword,
	}

	resp := db_ksc.Create(&post)

	if resp.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "Fail",
			"data":   post,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"data":   post,
		})
	}
}

func getPostByUrl(c *gin.Context) {
	url := c.Query("url")
	data := Posts{}
	db_ksc.Model(&Posts{}).Where("url=?", url).First(&data)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"data":   data,
	})
}

func updatePost(c *gin.Context) {
	var input Posts
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := db_ksc.Model(&input).Where("id = ?", input.ID).Updates(input)
	if res.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "Error!",
			"data":   input,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"data":   input,
		})
	}
}

func deletePost(c *gin.Context) {
	id := c.Query("id")
	uuid, _ := strconv.Atoi(id)
	data := Posts{}
	db_ksc.Model(&Posts{}).Where("id=?", uuid).First(&data)
	db_ksc.Delete(&data)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"data":   data,
	})
}
