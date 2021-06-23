package controllers

import (
	"net/http"

	"github.com/AlanHerediaG/test-jwt/auth"
	"github.com/AlanHerediaG/test-jwt/database"
	"github.com/AlanHerediaG/test-jwt/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = user.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = user.CreateUserRecord()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = user.CheckPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": err.Error(),
		})
		return
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	tokenResponse := LoginResponse{
		Token: signedToken,
	}
	c.JSON(http.StatusOK, tokenResponse)
}
