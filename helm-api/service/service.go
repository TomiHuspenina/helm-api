package service

import (
	"helm-api/model"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func GetImagesFromChartURL(url model.HelmRequest) ([]model.ImageInfo, error) {

	var settings = cli.New()
	chartDwldr := downloader.ChartDownloader{
		Out:     os.Stdout,
		Getters: getter.All(settings),
		Options: []getter.Option{},
	}

	tempDir := os.TempDir()                                          //temporal file
	file, _, err := chartDwldr.DownloadTo(url.ChartURL, tempDir, "") //download files
	if err != nil {
		return nil, err
	}

	chartReq, err := loader.Load(file) //load the file
	if err != nil {
		return nil, err
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), "default", "memory", func(format string, v ...interface{}) {})
	if err != nil {
		return nil, err
	}
	install := action.NewInstall(actionConfig)
	install.DryRun = true
	install.ClientOnly = true
	install.ReleaseName = "dummy"
	install.Namespace = "default"

	release, err := install.Run(chartReq, nil)
	if err != nil {
		return nil, err
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

/*
func ExtractImagesFromChart(chartURL string) ([]string, error) {
	// Crear una acción para template
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), "default", os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	install := action.NewInstall(actionConfig)
	install.DryRun = true
	install.ReleaseName = "dummy"
	install.Replace = true
	install.ClientOnly = true

	cp, err := install.ChartPathOptions.LocateChart(chartURL, settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	rel, err := install.Run(chart, nil)
	if err != nil {
		return nil, err
	}

	// Buscar imágenes en el YAML renderizado
	return extractImagesFromManifests(rel.Manifest), nil
}

func extractImagesFromManifests(manifest string) []string {
	var images []string
	docs := strings.Split(manifest, "\n---\n")

	for _, doc := range docs {
		var m map[string]interface{}
		err := yaml.Unmarshal([]byte(doc), &m)
		if err != nil || m == nil {
			continue
		}
		// Navegar hasta spec.template.spec.containers[].image
		spec, ok := m["spec"].(map[string]interface{})
		if !ok {
			continue
		}
		template, ok := spec["template"].(map[string]interface{})
		if !ok {
			continue
		}
		podSpec, ok := template["spec"].(map[string]interface{})
		if !ok {
			continue
		}
		containers, ok := podSpec["containers"].([]interface{})
		if !ok {
			continue
		}
		for _, c := range containers {
			container, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			if image, ok := container["image"].(string); ok {
				images = append(images, image)
			}
		}
	}
	return images
}*/
