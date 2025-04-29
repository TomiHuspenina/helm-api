package main

import (
	"fmt"
	"helm-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	router = gin.Default()
	//router.POST("/images", getImagesHandler)
	router.GET("/", func(c *gin.Context) {
		res, _ := service.GetImagesFromChartURL(c.Query("url"))
		c.JSON(200, res)
	})
	router.Run(":8080")

	fmt.Println("Server running in 8080")

}

func getImagesHandler(c *gin.Context) {

	var req string /*model.HelmRequest*/
	err := c.ShouldBindJSON(&req)
	if err != nil || req /*.ChartURL*/ == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid request body"})
		return
	}

	images, err := service.GetImagesFromChartURL(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, images)
	}
}
