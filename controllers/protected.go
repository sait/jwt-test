package controllers

import (
	"net/http"

	"github.com/AlanHerediaG/test-jwt/database"
	"github.com/AlanHerediaG/test-jwt/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(c *gin.Context) {
	var user models.User

	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": result.Error.Error(),
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"msg": result.Error.Error(),
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
