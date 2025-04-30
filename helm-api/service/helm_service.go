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

//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/mysql-10.1.0.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/apache-10.0.4.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/nginx-10.0.4.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/redis-17.0.0.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/mariadb-12.0.0.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/postgresql-12.0.0.tgz
//http://localhost:8080/?url=https://charts.bitnami.com/bitnami/wordpress-17.0.0.tgz

func GetImagesFromChartURL(url string) ([]model.ImageInfo, apierrors.ApiError) {

	var settings = cli.New()
	chartDwldr := downloader.ChartDownloader{ //descarga el chart (.tgz)
		Out:     os.Stdout,
		Getters: getter.All(settings),
		Options: []getter.Option{},
	}

	tempDir := os.TempDir()                                               //temporary file
	file, _, err := chartDwldr.DownloadTo(url /*.chartURL*/, tempDir, "") //download chart
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error downloading chart", err)
	}

	chartReq, err := loader.Load(file) //lee charthelm del archivo descargado para crear la estructura
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error loading chart", err)
	}

	actionConfig := new(action.Configuration) //entorno de instalacion falso
	err = actionConfig.Init(settings.RESTClientGetter(), "default", "memory", func(format string, v ...interface{}) {})
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error initializing action config", err)
	}
	install := action.NewInstall(actionConfig) //instalador de chart helm
	install.DryRun = true
	install.ClientOnly = true
	install.ReleaseName = "dummy"
	install.Namespace = "default"

	values := chartReq.Values

	release, err := install.Run(chartReq, values) //genera el contenido yaml del chart
	if err != nil {
		return nil, apierrors.NewInternalServerApiError("error running install dry-run", err)
	}

	var images []model.ImageInfo

	lines := strings.Split(release.Manifest, "\n") //separa el contenido en lineas
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

	layers, err := img.Layers() //extrae las capas de la imagen
	if err != nil {
		return 0, 0, err
	}

	var totalSize int64 = 0
	for _, layer := range layers {
		size, err := layer.Size() //obtiene el tamaño de cada capa
		if err != nil {
			return 0, 0, err
		}
		totalSize += size
	}

	return totalSize, len(layers), nil
}
