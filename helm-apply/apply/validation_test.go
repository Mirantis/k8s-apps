package apply

import (
	"reflect"
	"testing"
)

func TestCheckEmpty(t *testing.T) {
	cases := []struct {
		applier     Applier
		expectedErr error
	}{
		{
			Applier{
				map[string]string{"repo": ""},
				map[string]*Cluster{"cluster": nil},
				map[string]*Release{"release": nil},
			},
			nil,
		},
		{
			Applier{
				nil,
				map[string]*Cluster{"cluster": nil},
				map[string]*Release{"release": nil},
			},
			&ValidationError{"at least one repository under \"repos\" key should be defined"},
		},
		{
			Applier{
				map[string]string{"repo": ""},
				nil,
				map[string]*Release{"release": nil},
			},
			&ValidationError{"at least one cluster under \"clusters\" key should be defined"},
		},
		{
			Applier{
				map[string]string{"repo": ""},
				map[string]*Cluster{"cluster": nil},
				nil,
			},
			&ValidationError{"at least one release under \"releases\" key should be defined"},
		},
		{
			Applier{
				nil,
				nil,
				nil,
			},
			&ValidationError{"at least one repository under \"repos\" key should be defined"},
		},
	}
	for _, c := range cases {
		err := c.applier.checkEmpty()
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
		}
	}
}

func TestValidateRepositories(t *testing.T) {
	cases := []struct {
		applier     Applier
		expectedErr error
	}{
		{
			Applier{
				map[string]string{
					"mirantis": "https://mirantisworkloads.storage.googleapis.com",
					"local":    "http://127.0.0.1:8879"},
				nil,
				nil,
			},
			nil,
		},
		{
			Applier{
				map[string]string{
					"local": "127.0.0.1:8879"},
				nil,
				nil,
			},
			&ValidationError{"repo \"local\" has not valid URL: \"127.0.0.1:8879\""},
		},
		{
			Applier{
				map[string]string{
					"wrong": "/file/chart"},
				nil,
				nil,
			},
			&ValidationError{"repo \"wrong\" has not valid URL: \"/file/chart\""},
		},
	}
	for _, c := range cases {
		err := c.applier.validateRepositories()
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
		}
	}
}

func TestValidateClusters(t *testing.T) {
	cases := []struct {
		applier     Applier
		expectedErr error
	}{
		{
			Applier{
				nil,
				map[string]*Cluster{
					"first": {},
				},
				nil,
			},
			nil,
		},
		{
			Applier{
				nil,
				map[string]*Cluster{
					"first": {ExternalIP: "127.0.0.1"},
				},
				nil,
			},
			nil,
		},
		{
			Applier{
				nil,
				map[string]*Cluster{
					"first": {ExternalIP: "example.com"},
				},
				nil,
			},
			&ValidationError{"cluster \"first\" has not valid externalIP: \"example.com\""},
		},
	}
	for _, c := range cases {
		err := c.applier.validateClusters()
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
		}
	}
}

func TestValidateReleases(t *testing.T) {
	cases := []struct {
		applier     Applier
		expectedErr error
	}{
		{
			Applier{
				nil,
				nil,
				map[string]*Release{"test": {}},
			},
			&ValidationError{"\"cluster\" key should be defined for release \"test\""},
		},
		{
			Applier{
				nil,
				nil,
				map[string]*Release{"test": {ClusterName: "test-cluster"}},
			},
			&ValidationError{"\"chart\" key should be defined for release \"test\""},
		},
		{
			Applier{
				nil,
				nil,
				map[string]*Release{"test": {
					ClusterName: "test-cluster",
					ChartString: "test-repo/test-chart",
				}},
			},
			&ValidationError{"release \"test\" is referencing not defined cluster \"test-cluster\""},
		},
		{
			Applier{
				nil,
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{"test": {
					ClusterName: "test-cluster",
					ChartString: "test-repo/test-chart",
				}},
			},
			&ValidationError{"release \"test\" is referencing not defined chart repository \"test-repo\""},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{"test": {
					ClusterName: "test-cluster",
					ChartString: "test-repo/test-chart",
				}},
			},
			nil,
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{"test": {
					ClusterName: "test-cluster",
					ChartString: "test-repo/test-chart:1.2.0",
				}},
			},
			nil,
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{"test": {
					ClusterName: "test-cluster",
					ChartString: "test-repo",
				}},
			},
			&ValidationError{"test-repo: chart format is not valid. " +
				"should be \"<chart-repo-name>/<chart-name>:<chart-version>\" (chart version is optional)"},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{"test": {
					ClusterName:  "test-cluster",
					ChartString:  "test-repo/test-chart",
					Dependencies: map[string]string{"dep": "dep-release"},
				}},
			},
			&ValidationError{"release \"test\" depends on not defined release \"dep-release\""},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{
					"test": {
						ClusterName:  "test-cluster",
						ChartString:  "test-repo/test-chart",
						Dependencies: map[string]string{"dep": "dep-release"},
					},
					"dep-release": {
						ClusterName: "test-cluster",
						ChartString: "test-repo/dep-chart",
					},
				},
			},
			nil,
		},
	}
	for _, c := range cases {
		err := c.applier.validateReleases()
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
		}
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		applier     Applier
		expectedErr error
	}{
		{
			Applier{
				nil,
				nil,
				nil,
			},
			&ValidationError{"at least one repository under \"repos\" key should be defined"},
		},
		{
			Applier{
				map[string]string{"test-repo": "./repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{
					"test": {
						ClusterName: "test-cluster",
						ChartString: "test-repo/test-chart",
					},
				},
			},
			&ValidationError{"repo \"test-repo\" has not valid URL: \"./repo\""},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {
						ExternalIP: "test",
					},
				},
				map[string]*Release{
					"test": {
						ClusterName: "test-cluster",
						ChartString: "test-repo/test-chart",
					},
				},
			},
			&ValidationError{"cluster \"test-cluster\" has not valid externalIP: \"test\""},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{
					"test": {
						ClusterName:  "test-cluster",
						ChartString:  "test-repo/test-chart",
						Dependencies: map[string]string{"dep": "dep-release"},
					},
				},
			},
			&ValidationError{"release \"test\" depends on not defined release \"dep-release\""},
		},
		{
			Applier{
				map[string]string{"test-repo": "http://test.repo"},
				map[string]*Cluster{
					"test-cluster": {},
				},
				map[string]*Release{
					"test": {
						ClusterName:  "test-cluster",
						ChartString:  "test-repo/test-chart",
						Dependencies: map[string]string{"dep": "dep-release"},
					},
					"dep-release": {
						ClusterName: "test-cluster",
						ChartString: "test-repo/dep-chart",
					},
				},
			},
			nil,
		},
	}
	for _, c := range cases {
		err := c.applier.Validate()
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err: %q but actual is: %q", c.expectedErr, err)
		}
	}
}
