package client

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"helm-broker/utils"
	"k8s.io/helm/pkg/helm"
	"log"
)

// Install creates a new release
func Install(helmClient helm.Client, chartPath string, releaseName string, namespace string, values []byte) error {
	_, err := helmClient.ReleaseStatus(releaseName)
	if err == nil {
		log.Println("Release " + releaseName + " is already deployed")
		return nil
	}
	_, err = helmClient.InstallRelease(chartPath, namespace, helm.ReleaseName(releaseName),
		helm.ValueOverrides(values))
	return err
}

func GetConnectionStrings(helmClient helm.Client, releaseName string) (brokerapi.Credential, error) {
	status, err := helmClient.ReleaseStatus(releaseName)
	if err != nil {
		return brokerapi.Credential{}, err
	}
	return utils.GetConnectionStringsFromNotes(status.Info.Status.Notes)
}

// Delete deletes a release
func Delete(helmClient helm.Client, releaseName string) error {
	_, err := helmClient.DeleteRelease(releaseName)
	return err
}

func IsResourcesReady(helmClient helm.Client, releaseName string) (bool, error) {
	resources, err := getResources(helmClient, releaseName)
	if err != nil {
		return false, err
	}
	namespace, err := getNamespace(helmClient, releaseName)
	if err != nil {
		return false, err
	}
	return checkResourcesState(resources, namespace)
}

func GetSecrets(helmClient helm.Client, releaseName string) (brokerapi.Credential, error) {
	resources, err := getResources(helmClient, releaseName)
	if err != nil {
		return brokerapi.Credential{}, err
	}
	var secrets []string
	for typeRes, res := range resources {
		if typeRes == "Secret" {
			secrets = res
		}
	}
	namespace, err := getNamespace(helmClient, releaseName)
	if err != nil {
		return brokerapi.Credential{}, err
	}
	if len(secrets) != 0 {
		return getSecretsByNames(secrets, namespace)
	}
	return brokerapi.Credential{}, nil
}

func getResources(helmClient helm.Client, releaseName string) (map[string][]string, error) {
	status, err := helmClient.ReleaseStatus(releaseName)
	if err != nil {
		return map[string][]string{}, err
	}
	resourcesString := status.Info.Status.Resources
	return utils.ParseResources(resourcesString)
}

func getNamespace(helmClient helm.Client, releaseName string) (string, error) {
	status, err := helmClient.ReleaseStatus(releaseName)
	if err != nil {
		return "", err
	}
	return status.Namespace, nil
}
