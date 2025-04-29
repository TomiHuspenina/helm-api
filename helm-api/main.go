package main

import (
	"helm-api/app"
)

func main() {
	app.StartApp()
}

//func getImagesHandler(c *gin.Context) {
//
//	var req string /*model.HelmRequest*/
//	err := c.ShouldBindJSON(&req)
//	if err != nil || req /*.ChartURL*/ == "" {
//		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid request body"})
//		return
//	}
//
//	images, err := service.GetImagesFromChartURL(req)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//	} else {
//		c.JSON(http.StatusOK, images)
//	}
//}
