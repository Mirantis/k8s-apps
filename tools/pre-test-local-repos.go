package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	repoPathPtr = flag.String("repo", "charts/", "Path to charts repository")
	remoteRepo  = flag.String("remote-repo", "https://mirantisworkloads.storage.googleapis.com/", "Link to remote repo")
	revertRepo  = flag.Bool("revert", false, "Revert switching to local repo")
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func replaceRepo(path string) {
	reqFile := filepath.Join(path, "/requirements.yaml")
	yamlFile, err := ioutil.ReadFile(reqFile)
	if err != nil {
		log.Fatalf("Read file " + reqFile + " failed, not proceeding")
	}

	reqs := make(map[string][]map[string]string)
	yaml.Unmarshal(yamlFile, &reqs)
	deps := reqs["dependencies"]
	for _, dep := range deps {
		if *revertRepo {
			dep["repository"] = *remoteRepo
		} else {
			dep["repository"] = "file://../" + dep["name"] + "/"
		}
	}
	reqs["dependencies"] = deps

	b, _ := yaml.Marshal(reqs)
	ioutil.WriteFile(path+"/requirements.yaml", b, 0664)
}

func main() {
	flag.Parse()
	matches, err := filepath.Glob(*repoPathPtr + "/*")
	if err != nil {
		log.Fatalf("Failed to list directory %s: %s", *repoPathPtr, err)
	}
	var charts []string
	for _, match := range matches {
		charts = append(charts, filepath.Base(match))
	}

	for _, chart := range charts {
		chartPath := *repoPathPtr + chart
		stat, _ := os.Stat(chartPath)
		if stat.IsDir() {
			isExists, err := exists(chartPath + "/requirements.yaml")
			if err != nil {
				log.Fatalf("Requirements file searching gone wrong, not proceeding")
			}

			if isExists {
				replaceRepo(chartPath)
			}
		}
	}
}
