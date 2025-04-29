package app

import (
	"helm-api/service"

	"github.com/gin-gonic/gin"
)

func mapUrls() {
	//router.POST("/images", getImagesHandler)
	router.GET("/", func(c *gin.Context) {
		res, _ := service.GetImagesFromChartURL(c.Query("url"))
		c.JSON(200, res)
	})
}
