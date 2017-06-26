package utils

import (
	"errors"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"strings"
)

func GetConnectionStringsFromNotes(note string) (brokerapi.Credential, error) {
	underInternalURL := strings.Split(note, "Internal URL:\n")
	errMsg := "There is no connection strings"
	if len(underInternalURL) < 2 {
		return brokerapi.Credential{}, errors.New(errMsg)
	}
	strs := strings.Split(underInternalURL[1], "\n")
	if len(strs) == 0 {
		return brokerapi.Credential{}, errors.New(errMsg)
	}
	result := brokerapi.Credential{}

	for _, str := range strs {
		if str == "" {
			continue
		}
		subStrs := strings.SplitN(str, ":", 2)
		if len(subStrs) != 2 || subStrs[1] == "" {
			break
		}
		subStrs[0] = strings.TrimSpace(subStrs[0])
		subStrs[1] = strings.TrimSpace(subStrs[1])
		result[subStrs[0]] = subStrs[1]
	}
	if len(result) == 0 {
		return brokerapi.Credential{}, errors.New(errMsg)
	}
	return result, nil
}

func ParseResources(resources string) (map[string][]string, error) {
	strs := strings.Split(resources, "\n")
	res := map[string][]string{}
	key := ""
	for _, str := range strs {
		if strings.HasPrefix(str, "==>") {
			key = strings.Split(str, "/")[1]
		} else if strings.HasPrefix(str, "NAME") || str == "" {
			continue
		} else if key != "" {
			st := strings.Split(str, " ")
			res[key] = append(res[key], st[0])
		}
	}
	if len(res) == 0 {
		return res, errors.New("There is no resources")
	}
	return res, nil
}
