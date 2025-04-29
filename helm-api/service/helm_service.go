package service

import (
	apierrors "helm-api/error"
	"helm-api/model"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func GetImagesFromChartURL(url string) ([]model.ImageInfo, apierrors.ApiError) {

	var settings = cli.New()
	chartDwldr := downloader.ChartDownloader{
		Out:     os.Stdout,
		Getters: getter.All(settings),
		Options: []getter.Option{},
	}

	tempDir := os.TempDir()                                               //temporary file
	file, _, err := chartDwldr.DownloadTo(url /*.chartURL*/, tempDir, "") //download chart
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error downloading chart", err)
	}

	chartReq, err := loader.Load(file) //load the file
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error loading chart", err)
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), "default", "memory", func(format string, v ...interface{}) {})
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error initializing action config", err)
	}
	install := action.NewInstall(actionConfig)
	install.DryRun = true
	install.ClientOnly = true
	install.ReleaseName = "dummy"
	install.Namespace = "default"

	values := chartReq.Values

	release, err := install.Run(chartReq, values)
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error running install dry-run", err)
	}

	var images []model.ImageInfo

	lines := strings.Split(release.Manifest, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "image:") {
			image := strings.TrimSpace(strings.TrimPrefix(line, "image:"))

			sizeB, layers, err := getImagesDetail(image)
			if err != nil {
				images = append(images, model.ImageInfo{
					Image:     image,
					Size:      0,
					NumLayers: 0,
				})
				continue
			}

			sizeMB := float64(sizeB) / 1024 / 1024

			images = append(images, model.ImageInfo{
				Image:     image,
				Size:      sizeMB,
				NumLayers: layers,
			})
		}
	}

	return images, nil
}

func getImagesDetail(imageRef string) (int64, int, error) {

	img, err := crane.Pull(imageRef)
	if err != nil {
		return 0, 0, err
	}

	layers, err := img.Layers()
	if err != nil {
		return 0, 0, nil
	}

	var totalSize int64 = 0
	for _, layer := range layers {
		size, err := layer.Size()
		if err != nil {
			return 0, 0, err
		}
		totalSize += size
	}

	return totalSize, len(layers), nil
}
