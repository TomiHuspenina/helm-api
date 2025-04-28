package main

import (
	"encoding/json"
	"fmt"
	"helm-api/model"
	"helm-api/service"
	"net/http"
)

func main() {
	http.HandleFunc("/get-images", getImagesHandler)

	fmt.Println("Servidor corriendo en :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getImagesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req model.HelmRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ChartURL == "" {
		http.Error(w, fmt.Sprintf("error decoding json", err), http.StatusBadRequest)
		return
	}

	images, err := service.GetImagesFromChartURL(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting images from url", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}
