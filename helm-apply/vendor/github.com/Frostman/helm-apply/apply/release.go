package apply

import (
	"fmt"

	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/Frostman/helm-apply/getter"
	"github.com/databus23/helm-diff/manifest"
	"github.com/imdario/mergo"

	"github.com/Frostman/helm-apply/diff"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	hapi "k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/timeconv"
)

type Release struct {
	Name              string              `yaml:"-"`
	ChartString       string              `yaml:"chart"`
	Chart             *Chart              `yaml:"-"`
	ClusterName       string              `yaml:"cluster"`
	Cluster           *Cluster            `yaml:"-"`
	Namespace         string              `yaml:"namespace"`
	Parameters        chartutil.Values    `yaml:"parameters"`
	Values            chartutil.Values    `yaml:"-"`
	Dependencies      map[string]string   `yaml:"dependencies"`
	DepReleases       map[string]*Release `yaml:"-"`
	InternalAddresses map[string]string   `yaml:"-"`
	ExternalAddresses map[string]string   `yaml:"-"`
	WaitFlag          bool                `yaml:"wait"`
	currRelease       *hapi.Release
}

type Chart struct {
	Repository string
	Name       string
	Version    string
}

func NewChart(repository, name, version string) *Chart {
	return &Chart{repository, name, version}
}

func (r *Release) ParseValues() ([]byte, error) {
	values, err := yaml.Marshal(r.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %s", err)
	}

	return values, nil
}

func PrintStatus(out io.Writer, res *services.GetReleaseStatusResponse) {
	if res.Info.LastDeployed != nil {
		fmt.Fprintf(out, "LAST DEPLOYED: %s\n", timeconv.String(res.Info.LastDeployed))
	}
	fmt.Fprintf(out, "NAMESPACE: %s\n", res.Namespace)
	fmt.Fprintf(out, "STATUS: %s\n", res.Info.Status.Code)
	fmt.Fprintf(out, "\n")
	if len(res.Info.Status.Resources) > 0 {
		re := regexp.MustCompile("  +")

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "RESOURCES:\n%s\n", re.ReplaceAllString(res.Info.Status.Resources, "\t"))
		w.Flush()
	}

	if len(res.Info.Status.Notes) > 0 {
		fmt.Fprintf(out, "NOTES:\n%s\n", res.Info.Status.Notes)
	}
}

func (r *Release) parseInternalAddresses(notes string) (map[string]string, error) {
	underInternalURL := strings.Split(notes, "Internal URL:\n")
	if len(underInternalURL) < 2 {
		return nil, fmt.Errorf("no \"Internal URL:\" section was found")
	}
	strs := strings.Split(underInternalURL[1], "\n")
	errorMsg := "there are no <service>:<addr> pairs under \"Internal URL:\" section"
	if len(strs) == 0 {
		return nil, fmt.Errorf(errorMsg)
	}
	addresses := map[string]string{}
	for _, str := range strs {
		if str == "" {
			continue
		}
		subStrs := strings.SplitN(str, ":", 2)
		if len(subStrs) != 2 || subStrs[1] == "" {
			break
		}
		addresses[strings.TrimSpace(subStrs[0])] = strings.TrimSpace(subStrs[1])
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf(errorMsg)
	}
	return addresses, nil
}

func (r *Release) parseExternalAddresses(notes string) (map[string]string, error) {
	extServices := strings.Split(notes, "External services:\n")
	if len(extServices) < 2 {
		return nil, fmt.Errorf("no \"External services:\" section was found")
	}
	strs := strings.Split(extServices[1], "\n")
	errorMsg := "there are no <service>: <service_name>:<port> under \"External services:\" section"
	if len(strs) == 0 {
		return nil, fmt.Errorf(errorMsg)
	}
	addresses := map[string]string{}
	for _, str := range strs {
		if str == "" {
			continue
		}
		subStrs := strings.SplitN(str, ":", 3)
		if len(subStrs) != 3 {
			break
		}
		var extAddr string
		if dryRunFlag {
			extAddr = "not-defined-during-dry-run:80"
		} else {
			var err error
			extAddr, err = r.Cluster.serviceExternalAddress(
				strings.TrimSpace(subStrs[1]), strings.TrimSpace(subStrs[2]), r.GetNamespace())
			if err != nil {
				fmt.Println("error getting external address: ", err)
			}
		}
		if extAddr != "" {
			addresses[strings.TrimSpace(subStrs[0])] = extAddr
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf(errorMsg)
	}
	return addresses, nil
}

func (r *Release) injectDependencies() error {
	for chartName, release := range r.DepReleases {
		var addresses map[string]string
		if r.ClusterName == release.ClusterName {
			addresses = release.InternalAddresses
		} else {
			addresses = release.ExternalAddresses
		}
		if len(addresses) == 0 {
			return fmt.Errorf("release \"%s\" doesn't have addresses to refer to", release.Name)
		}
		values := chartutil.Values{}
		for k, v := range release.Values {
			values[k] = v
		}
		values["deployChart"] = false
		values["addresses"] = addresses
		depParameters := chartutil.Values{chartName: values}
		if r.Parameters == nil {
			r.Parameters = chartutil.Values{}
		}
		if err := mergo.Merge(&r.Parameters, depParameters); err != nil {
			return err
		}
	}
	return nil
}

func (r Release) GetNamespace() (namespace string) {
	namespace = r.Namespace
	if namespace == "" {
		namespace = r.Cluster.GetDefaultNamespace()
	}
	return
}

func (r Release) Download() (chartPath string, err error) {
	chartUrl, err := repo.FindChartInRepoURL(r.Chart.Repository, r.Chart.Name, r.Chart.Version,
		"", "", "", getterProviders)

	if err != nil {
		return "", err
	}
	chartDst, err := ioutil.TempFile("", r.Chart.Name)
	if err != nil {
		return "", fmt.Errorf("cannot write index file for repository requested")
	}
	chartGetter, err := getter.NewHTTPGetter(chartUrl, "", "", "")
	if err != nil {
		panic(err)
	}
	resp, err := chartGetter.Get(chartUrl)
	if err != nil {
		return "", err
	}
	_, err = chartDst.Write(resp.Bytes())
	if err != nil {
		return "", err
	}
	return chartDst.Name(), nil
}

func (r *Release) Install() error {
	fmt.Printf("Processing \"%s\" release...\n", r.Name)
	chartDst, err := r.Download()
	if err != nil {
		return fmt.Errorf("cannot download chart: %s", err)
	}
	defer os.Remove(chartDst)

	err = r.injectDependencies()
	if err != nil {
		return err
	}

	namespace := r.GetNamespace()
	values, err := r.ParseValues()
	if err != nil {
		return err
	}
	helmClient, err := r.Cluster.GetHelmClient()
	if err != nil {
		return err
	}
	prevReleaseResp, err := helmClient.ReleaseContent(r.Name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		fmt.Printf("Release %q does not exist. Installing it now...\n", r.Name)
		currReleaseResp, err := helmClient.InstallRelease(chartDst, namespace, helm.ReleaseName(r.Name),
			helm.ValueOverrides(values), helm.InstallDryRun(dryRunFlag))
		if err != nil {
			return err
		}
		r.currRelease = currReleaseResp.Release
		if diffFlag {
			currManifest := manifest.Parse(currReleaseResp.Release.Manifest)
			fmt.Printf("Diff for \"%s\" release:\n", r.Name)
			diff.Manifests(map[string]string{}, currManifest, os.Stdout)
		}
	} else {
		fmt.Printf("Release %q already exists. Updating it now...\n", r.Name)
		currReleaseResp, err := helmClient.UpdateRelease(r.Name, chartDst, helm.UpdateValueOverrides(values),
			helm.UpgradeDryRun(dryRunFlag))
		if err != nil {
			return err
		}
		r.currRelease = currReleaseResp.Release
		if diffFlag {
			prevManifest := manifest.Parse(prevReleaseResp.Release.Manifest)
			currManifest := manifest.Parse(currReleaseResp.Release.Manifest)
			fmt.Printf("Diff for \"%s\" release:\n", r.Name)
			diff.Manifests(prevManifest, currManifest, os.Stdout)
		}
	}

	if !dryRunFlag {
		status, err := helmClient.ReleaseStatus(r.Name)
		if err != nil {
			return err
		}

		if verboseFlag {
			PrintStatus(os.Stdout, status)
		}
	}
	if intAddresses, err := r.parseInternalAddresses(r.currRelease.Info.Status.Notes); err == nil {
		r.InternalAddresses = intAddresses
	}
	if extAddresses, err := r.parseExternalAddresses(r.currRelease.Info.Status.Notes); err == nil {
		r.ExternalAddresses = extAddresses
	}

	mergedValues, err := chartutil.CoalesceValues(r.currRelease.Chart, r.currRelease.Config)
	if err != nil {
		return err
	}
	r.Values = mergedValues

	return nil
}

func (r *Release) Wait() error {
	if r.WaitFlag && !dryRunFlag {
		fmt.Printf("Waiting for \"%s\" release deployment...\n", r.Name)
		helmClient, err := r.Cluster.GetHelmClient()
		if err != nil {
			return err
		}
		_, err = helmClient.UpdateReleaseFromChart(
			r.Name, r.currRelease.Chart, helm.UpdateValueOverrides([]byte{}), helm.UpgradeWait(true))
		if err != nil {
			return err
		}
		fmt.Printf("\"%s\" release has been deployed.\n", r.Name)
	}
	return nil
}

func (r *Release) PrintAddresses() {
	if len(r.ExternalAddresses) > 0 {
		fmt.Printf("\"%s\" external addresses:\n", r.Name)
		for k, v := range r.ExternalAddresses {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
}
