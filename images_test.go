package images_test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestImages(t *testing.T) {
	versions := ImageVersions(t)
	matches, err := filepath.Glob("charts/*")
	if err != nil {
		t.Fatalf("Failed to list charts directory: %s", err)
	}
	for _, match := range matches {
		chart := path.Base(match)
		t.Logf("Check %s chart", chart)
		t.Run(chart, func(t *testing.T) {
			values, err := ioutil.ReadFile(path.Join(match, "values.yaml"))
			if err != nil {
				t.Fatalf("Failed to read values.yaml for %s chart", chart)
			}
			v := make(map[interface{}]interface{})
			err = yaml.Unmarshal(values, &v)
			if err != nil {
				t.Fatalf("Error parsing values.yaml file")
			}
			image := v["image"]
			found := false
			if image != nil {
				CheckVersion(t, image, versions)
				found = true
			}
			for key, item := range v {
				value, ok := item.(map[interface{}]interface{})
				if ok {
					image = value["image"]
					if image != nil {
						t.Logf("Found image in %s section", key)
						CheckVersion(t, image, versions)
						found = true
					}
				}
			}
			if !found {
				t.Logf("Image section not found for %s chart", chart)
			}
		})
	}
}

func ImageVersions(t *testing.T) map[string]string {
	versions := make(map[string]string)
	matches, err := filepath.Glob("images/*")
	if err != nil {
		t.Fatalf("Failed to list images/ directory: %s", err)
	}
	for _, match := range matches {
		imageName := filepath.Base(match)
		version, err := ioutil.ReadFile(path.Join(match, ".version"))
		if err != nil {
			t.Fatalf("Couldn't get %s image version: %s", imageName, err)
		}
		versions[imageName] = strings.TrimSpace(string(version))
	}
	return versions
}

func CheckVersion(t *testing.T, image interface{}, versions map[string]string) {
	img, ok := image.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Incorrect image section\n%+v", image)
	}
	if img["repository"].(string) == "mirantisworkloads/" {
		imgName := img["name"].(string)
		imgVersion := img["tag"].(string)
		foundVersion := versions[imgName]
		if foundVersion == "" {
			t.Fatalf("Image not found %s", imgName)
		}
		if foundVersion != imgVersion {
			t.Fatalf("Image version in chart %s, but found %s", imgVersion, foundVersion)
		}
	} else {
		t.Logf("Found trird-party image, skipping: %+v", img)
	}
}
