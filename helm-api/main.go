package main

import (
	"fmt"
	"helm-api/model"
	"helm-api/service"
)

func main() {
	req := model.HelmRequest{
		ChartURL: "https://charts.bitnami.com/bitnami/nginx-15.0.0.tgz",
	}

	images, err := service.GetImagesFromChartURL(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, img := range images {
		fmt.Printf("Image: %s, Size: %s, Layers: %d\n", img.Image, img.Size, img.NumLayers)
	}
}
