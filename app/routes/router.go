package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", handlePong)
	return router
}

func handlePong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
