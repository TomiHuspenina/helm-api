package main

import (
	"fmt"
	"helm-api/model"
	"helm-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	router = gin.Default()
	router.POST("/images", getImagesHandler)
	router.Run(":8080")

	fmt.Println("Servidor corriendo en :8080...")

}

func getImagesHandler(c *gin.Context) {

	var req model.HelmRequest
	err := c.ShouldBindJSON(&req)
	if err != nil || req.ChartURL == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid request body"})
		return
	}

	images, err := service.GetImagesFromChartURL(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": ""})
	}

	c.JSON(http.StatusOK, images)
}
