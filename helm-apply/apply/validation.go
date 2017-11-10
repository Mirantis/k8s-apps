package apply

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (applier *Applier) Validate() error {
	if err := applier.checkEmpty(); err != nil {
		return err
	}

	if err := applier.validateRepositories(); err != nil {
		return err
	}

	if err := applier.validateClusters(); err != nil {
		return err
	}

	if err := applier.validateReleases(); err != nil {
		return err
	}

	return nil
}

func (applier *Applier) checkEmpty() error {
	if len(applier.Repositories) == 0 {
		return &ValidationError{"at least one repository under \"repos\" key should be defined"}
	}
	if len(applier.Clusters) == 0 {
		return &ValidationError{"at least one cluster under \"clusters\" key should be defined"}
	}
	if len(applier.Releases) == 0 {
		return &ValidationError{"at least one release under \"releases\" key should be defined"}
	}
	return nil
}

func (applier *Applier) validateRepositories() error {
	for name, url := range applier.Repositories {
		if !govalidator.IsURL(url) {
			return &ValidationError{fmt.Sprintf(
				"repo \"%s\" has not valid URL: \"%s\"", name, url)}
		}
	}
	return nil
}

func (applier *Applier) validateClusters() error {
	for name, cluster := range applier.Clusters {
		if len(cluster.ExternalIP) != 0 && !govalidator.IsIP(cluster.ExternalIP) {
			return &ValidationError{fmt.Sprintf(
				"cluster \"%s\" has not valid externalIP: \"%s\"", name, cluster.ExternalIP)}
		}
	}
	return nil
}

func (applier *Applier) validateReleases() error {
	for releaseName, release := range applier.Releases {
		if len(release.ClusterName) == 0 {
			return &ValidationError{fmt.Sprintf(
				"\"cluster\" key should be defined for release \"%s\"", releaseName)}
		}
		if len(release.ChartString) == 0 {
			return &ValidationError{fmt.Sprintf(
				"\"chart\" key should be defined for release \"%s\"", releaseName)}
		}
		if _, ok := applier.Clusters[release.ClusterName]; !ok {
			return &ValidationError{fmt.Sprintf(
				"release \"%s\" is referencing not defined cluster \"%s\"", releaseName, release.ClusterName)}
		}
		repo, _, _, err := parseChartString(release.ChartString)
		if err != nil {
			return &ValidationError{err.Error()}
		}
		if _, ok := applier.Repositories[repo]; !ok {
			return &ValidationError{fmt.Sprintf(
				"release \"%s\" is referencing not defined chart repository \"%s\"", releaseName, repo)}
		}

		for _, depReleaseName := range release.Dependencies {
			if _, ok := applier.Releases[depReleaseName]; !ok {
				return &ValidationError{fmt.Sprintf(
					"release \"%s\" depends on not defined release \"%s\"", releaseName, depReleaseName)}
			}
		}
	}
	return nil
}
