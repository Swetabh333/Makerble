package middleware

import (
	"log"
	"net/http"

	"github.com/Swetabh333/Makerble/app/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// This middleware is used to protect our routes by checking if the requesting user is already logged in via a jwt token , if the token is expired but the refresh token is valid then the user is provided a new token

func VerifyAuthentication(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		// If the token is invalid, redirect to login

		accessToken, err := helper.ValidateToken(token)
		if err == nil && accessToken.Valid {
			claims := accessToken.Claims.(jwt.MapClaims)
			userID, ok := claims["sub"].(string)
			if !ok {
				log.Println("ID not found")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized to access this resource, please login to continue",
				})
				c.Abort()
				return
			}
			role, ok := claims["role"].(string)
			if !ok {
				log.Println("role not found")

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized to access this resource, please login to continue",
				})
				c.Abort()
				return
			}
			c.Set("ID", uuid.MustParse(userID))
			c.Set("Role", role)
			c.Next()
		}
	}
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		log.Println("Refresh token not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized to access this resource, please login to continue",
		})

		c.Abort()
		return
	}

	refreshTokenObj, err := helper.ValidateRefreshToken(refreshToken)
	if err != nil || !refreshTokenObj.Valid {
		// If the refresh token is invalid, redirect to login
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized to access this resource, please login to continue",
		})

		c.Abort()
		return
	}

	claims := refreshTokenObj.Claims.(jwt.MapClaims)
	newAccessToken, err := helper.GenerateToken(claims["user"].(string), claims["sub"].(string), claims["role"].(string))
	if err != nil {
		// Handle error if new token generation fails
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized to access this resource, please login to continue",
		})
		c.Abort()
		return
	}
	c.SetCookie("token", newAccessToken, 24*3600, "/", "", false, true)
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized to access this resource, please login to continue",
		})
		c.Abort()
		return
	}
	role, ok := claims["role"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized to access this resource, please login to continue",
		})

		c.Abort()
		return
	}
	c.Set("ID", uuid.MustParse(userID))
	c.Set("Role", role)

	c.Next()
}
