package service

import (
	apierrors "helm-api/error"
	"helm-api/model"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func GetImagesFromChartURL(url model.HelmRequest) ([]model.ImageInfo, apierrors.ApiError) {

	var settings = cli.New()
	chartDwldr := downloader.ChartDownloader{
		Out:     os.Stdout,
		Getters: getter.All(settings),
		Options: []getter.Option{},
	}

	tempDir := os.TempDir()                                          //temporary file
	file, _, err := chartDwldr.DownloadTo(url.ChartURL, tempDir, "") //download chart
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

	release, err := install.Run(chartReq, nil)
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error running install dry-run", err)
	}

	var images []model.ImageInfo

	lines := strings.Split(release.Manifest, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "image:") {
			image := strings.TrimSpace(strings.TrimPrefix(line, "image:"))

			images = append(images, model.ImageInfo{
				Image:     image,
				Size:      "",
				NumLayers: 0,
			})
		}
	}

	return images, nil
}
