package model

type HelmRequest struct {
	ChartURL string `json:"chart_url" binding:"required,url"`
}

type ImageInfo struct {
	Image     string `json:"image"`
	Size      string `json:"size"`
	NumLayers int    `json:"num_layers"`
}

type HelmResponse struct {
	Images []ImageInfo `json:"images"`
}
