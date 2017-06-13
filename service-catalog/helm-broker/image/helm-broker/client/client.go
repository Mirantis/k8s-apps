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
	conn_strings, err := utils.GetConnectionStringsFromNotes(status.Info.Status.Notes)
	return conn_strings, err
}

// Delete deletes a release
func Delete(helmClient helm.Client, releaseName string) error {
	_, err := helmClient.DeleteRelease(releaseName)
	return err
}
