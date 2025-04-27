package main

import (
	"Obmondo/model"
	"Obmondo/service"
	"fmt"
)

func main() {
	req := model.HelmRequest{
		ChartURL: "https://charts.bitnami.com/bitnami/nginx-15.0.0.tgz", // ejemplo real
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
