package controller

import (
	"errors"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"gopkg.in/yaml.v2"
)

func validateInstanceRequest(req *brokerapi.CreateServiceInstanceRequest) error {
	name, isExist := req.Parameters["name"]
	if !isExist {
		return errors.New("There is no name in the instance creation request")
	}
	_, isSuccess := name.(string)
	if !isSuccess {
		return errors.New("Can't convert to string the name from the instance creation request")
	}

	version, isExist := req.Parameters["version"]
	if !isExist {
		return errors.New("There is no version in the instance creation request")
	}
	_, isSuccess = version.(string)
	if !isSuccess {
		return errors.New("Can't convert to string the version from the instance creation request")
	}

	namespace, isExist := req.Parameters["namespace"]
	if !isExist {
		return errors.New("There is no namespace in the instance creation request")
	}
	_, isSuccess = namespace.(string)
	if !isSuccess {
		return errors.New("Can't convert to string the namespace from the instance creation request")
	}

	values, isExist := req.Parameters["namespace"]
	if !isExist {
		return errors.New("There is no namespace in the instance creation request")
	}
	_, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	return nil
}
