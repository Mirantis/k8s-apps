package apply

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/getter"
)

var (
	verboseFlag     bool
	diffFlag        bool
	dryRunFlag      bool
	getterProviders getter.Providers
)

type Applier struct {
	Repositories map[string]string   `yaml:"repos"`
	Clusters     map[string]*Cluster `yaml:"clusters"`
	Releases     map[string]*Release `yaml:"releases"`
}

func NewApplier(configFilePath string, providers getter.Providers) (*Applier, error) {
	applier := Applier{}
	err := applier.parseConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	err = applier.Validate()
	if err != nil {
		return nil, err
	}
	getterProviders = providers
	err = applier.prepareApplier()
	if err != nil {
		return nil, err
	}
	return &applier, nil
}

func (applier *Applier) Run(verbose, diff, dryRun bool) error {
	verboseFlag = verbose
	diffFlag = diff
	dryRunFlag = dryRun

	installOrder, err := ResolveDependencies(applier.Releases)
	if err != nil {
		return err
	}
	for _, release := range installOrder {
		err := applier.Releases[release].Install()
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n")
	for _, release := range installOrder {
		err := applier.Releases[release].Wait()
		if err != nil {
			return err
		}
	}
	fmt.Printf("\n")
	for _, release := range installOrder {
		applier.Releases[release].PrintAddresses()
	}
	return nil
}

func (applier *Applier) parseConfig(configFilePath string) error {
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("could not open config file %s: %s", configFilePath, err)
	}

	err = yaml.Unmarshal(configFile, &applier)
	if err != nil {
		return fmt.Errorf("failed to parse config file %s: %s", configFilePath, err)
	}
	return nil
}

func (applier *Applier) prepareApplier() error {
	for name, cluster := range applier.Clusters {
		cluster.Name = name
	}

	for name, release := range applier.Releases {
		release.Name = name

		repository, name, version, err := parseChartString(release.ChartString)
		if err != nil {
			return err
		}
		release.Chart = NewChart(applier.Repositories[repository], name, version)
		release.Cluster = applier.Clusters[release.ClusterName]
	}

	return nil
}

func parseChartString(chartString string) (repository, name, version string, err error) {
	myExp := regexp.MustCompile(`^(?P<repository>[[:alnum:]-_]+)/(?P<name>[[:alnum:]-_]+)(:(?P<version>.+))?$`)
	match := myExp.FindStringSubmatch(chartString)
	if len(match) == 0 {
		return "", "", "", fmt.Errorf("%s: chart format is not valid. should be "+
			"\"<chart-repo-name>/<chart-name>:<chart-version>\" (chart version is optional)", chartString)
	}
	repository = match[1]
	name = match[2]

	if len(match) == 5 {
		version = match[4]
	}
	return
}
