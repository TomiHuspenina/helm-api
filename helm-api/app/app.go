package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func StartApp() {
	router = gin.Default()
	mapUrls()
	router.Run(":8080")
	fmt.Println("Server running in 8080")
}
