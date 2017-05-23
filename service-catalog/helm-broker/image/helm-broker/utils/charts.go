package utils

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
)

const (
	ChartsPath = "/opt/helm-broker/charts/"
	IndexName  = "index.yaml"
)

type Index struct {
	ApiVersion string             `yaml:"apiVersion"`
	Entries    map[string][]Chart `yaml:"entries"`
	Generated  string             `yaml:"generated"`
}

type Chart struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Urls        []string `yaml:"urls"`
	Description string   `yaml:"description"`
}

func DownloadIndex(urlCharts string) error {
	u, err := url.Parse(urlCharts)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, IndexName)
	indexUrl := u.String()
	return downloadFile(ChartsPath, IndexName, indexUrl)
}

func ParseIndex() (Index, error) {
	indexPath := path.Join(ChartsPath, IndexName)
	yamlFile, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return Index{}, err
	}
	var index Index
	err = yaml.Unmarshal(yamlFile, &index)

	return index, err
}

func DownloadChart(name, version string) (string, error) {
	index, err := ParseIndex()
	if err != nil {
		return "", err
	}
	chartInfo, isExist := index.Entries[name]
	if !isExist {
		return "", errors.New("There is no a chart with the name " + name)
	}
	var chartUrl string
	for _, info := range chartInfo {
		if info.Version == version {
			chartUrl = info.Urls[0]
			break
		}
	}
	if chartUrl == "" {
		return "", errors.New("There is no URL for chart " + name + version + " in the index.yaml")
	}
	ls := strings.Split(chartUrl, "/")
	tarName := ls[len(ls)-1]
	chartPath := path.Join(ChartsPath, tarName)
	err = downloadFile(ChartsPath, tarName, chartUrl)
	return chartPath, err
}
