package go_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

func LookupEnvDefault(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	} else {
		return def
	}
}

var (
	repoPathPtr        = flag.String("repo", "charts/", "Path to charts repository")
	imagesPathPtr      = flag.String("images-dir", "images/", "Path to Dockerfiles")
	configPathPtr      = flag.String("config", "tests/", "Path to charts config files")
	excludePtr         = flag.String("exclude", "", "List of charts to exclude from run")
	prefixPtr          = flag.String("prefix", "", "Prefix to prepend to object names (releases, namespaces)")
	chartsPtr          = flag.Bool("charts", true, "Test charts")
	imagesPtr          = flag.Bool("images", false, "Build images")
	pushPtr            = flag.Bool("push", false, "Push images after tests")
	imageRepoPtr       = flag.String("image-repo", "127.0.0.1:5000", "Image repo address for test")
	pushRepoPtr        = flag.String("push-repo", "mirantisworkloads", "Image repo address for push")
	buildImagesOptsPtr = flag.String("build-images-opts", "", "Docker opts for building images")
	verifyVersion = flag.Bool("verify-version", false, "Run tests to verify new helm/k8s version")
	remoteCluster = flag.String("remote-cluster", "127.0.0.1", "Cluster IP if remote kubernetes cluster is used")
	verifyIngress = flag.Bool("verify-ingress", false, "Ensure ingress is working correctly")
	ingressSvc    = flag.String("ingress-svc", "", "Ingress service host:port to connect the ingress resource")
	helmCmd            = LookupEnvDefault("HELM_CMD", "helm")
	kubectlCmd         = LookupEnvDefault("KUBECTL_CMD", "kubectl")
)

func TestVerify(t *testing.T) {
	if *verifyVersion {
		t.Run("verify_version", VerifyVersion)
	}
	if *verifyIngress {
		t.Run("verify_ingress", VerifyIngress)
	}
}

func TestRoot(t *testing.T) {
	if *imagesPtr {
		t.Run("images", RunImages)
	}
	if *chartsPtr {
		t.Run("charts", RunCharts)
	}
	if *pushPtr {
		t.Run("push", RunPushImages)
	}
}

func VerifyVersion(t *testing.T) {
	// 10.3 verify
	repoAddResult := RunCmdTest(t, "repo_add", helmCmd, "repo", "add", "mirantisworkloads", "https://mirantisworkloads.storage.googleapis.com")
	if !repoAddResult {
		t.Fatalf("Adding repo failed, not proceeding")
	}

	kafkaChart := path.Join(*repoPathPtr, "/kafka")

	// 4, 10.6 verify
	depUpArgs := []string{helmCmd, "dep", "up", kafkaChart}
	depUpResult := RunCmdTest(t, "dep_up", depUpArgs...)
	if !depUpResult {
		t.Fatalf("Dependencies update failed for kafka, not proceeding")
	}

	// 10.2 verify
	lintResult := RunCmdTest(t, "lint", helmCmd, "lint", kafkaChart)
	if !lintResult {
		t.Fatalf("lint failed, not proceeding")
	}

	// 1, 2, 4, 10.1 verify
	installKafkaArgs := "persistence.type=PersistentVolumeClaim"
	kafkaNs := *prefixPtr + randStringRunes(10)
	kafkaRel := *prefixPtr + randStringRunes(8)
	t.Run("create_chart", func(t *testing.T) {
		CreateChartVersionVerify(t, kafkaChart, installKafkaArgs, "", kafkaNs, kafkaRel)
	})

	// 10.4 verify
	listResult := RunCmdTest(t, "list", helmCmd, "list")
	if !listResult {
		t.Fatalf("list releases failed, not proceeding")
	}

	// 7 verify
	t.Run("delete_chart", func(t *testing.T) {
		DeleteChartVersionVerify(t, kafkaNs, kafkaRel)
	})

	escChart := path.Join(*repoPathPtr, "/elasticsearch")

	// 3 verify
	installEscArgs := "client.service.type=NodePort,client.service.nodePort=31111"
	escNs := *prefixPtr + randStringRunes(10)
	escRel := *prefixPtr + randStringRunes(3)
	t.Run("create_chart", func(t *testing.T) {
		CreateChartVersionVerify(t, escChart, installEscArgs, "", escNs, escRel)
	})

	_, nodePortErr := http.Get("http://" + *remoteCluster + ":31111/")
	if nodePortErr != nil {
		t.Fatalf("NodePort works incorrectly, not proceeding")
	}
}

func VerifyIngress(t *testing.T) {
	// 5 verify
	hdfsChart := path.Join(*repoPathPtr, "/hdfs")
	hdfsNs := *prefixPtr + randStringRunes(10)
	hdfsRel := *prefixPtr + randStringRunes(8)

	t.Run("create_chart", func(t *testing.T) {
		CreateChartVersionVerify(t, hdfsChart, "", "tests/hdfs/ingress.yaml", hdfsNs, hdfsRel)
	})

	t.Run("check_ingress", func(t *testing.T) {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", *ingressSvc, nil)
		req.Header.Add("Host", "hdfs.ingress")
		client.Do(req)
	})

	t.Run("delete_chart", func(t *testing.T) {
		DeleteChartVersionVerify(t, hdfsNs, hdfsRel)
	})
}

func DiscoverArtifacts(t *testing.T, path string) []string {
	matches, err := filepath.Glob(path + "/*")
	if err != nil {
		t.Fatalf("Failed to list directory %s: %s", path, err)
	}
	var allArtifacts []string
	for _, match := range matches {
		allArtifacts = append(allArtifacts, filepath.Base(match))
	}
	t.Logf("Found artifacts: %+v", allArtifacts)
	var artifacts []string
	if flag.NArg() == 0 {
		artifacts = allArtifacts
	} else {
		matchSet := make(map[string]bool)
		for _, match := range matches {
			matchSet[filepath.Base(match)] = true
		}
		var notFound []string
		for _, arg := range flag.Args() {
			artifact := filepath.Base(arg)
			if !matchSet[artifact] {
				notFound = append(notFound, arg)
			} else {
				artifacts = append(artifacts, artifact)
			}
		}
		if len(notFound) > 0 {
			t.Fatalf("Couldn't find these artifacts: %+v", notFound)
		}
	}
	return artifacts
}

func ListExcludesCharts() map[string]bool {
	excludes := make(map[string]bool)
	for _, e := range strings.Split(*excludePtr, ",") {
		excludes[filepath.Base(e)] = true
	}
	return excludes
}

func RunImages(t *testing.T) {
	for _, image := range DiscoverArtifacts(t, "images") {
		image := image
		t.Run(image, func(t *testing.T) {
			t.Parallel()
			RunImage(t, image)
		})
	}
}

func RunImage(t *testing.T, image string) {
	imageDir := path.Join(*imagesPathPtr, image)
	imageVersionFilePath := path.Join(imageDir, ".version")
	version, err := ioutil.ReadFile(imageVersionFilePath)
	if err != nil {
		t.Fatalf("Couldn't get image version from file: %s\n%s", imageVersionFilePath, err)
	}
	imageTag := fmt.Sprintf("%s/%s:%s", *imageRepoPtr, image, strings.TrimSpace(string(version)))
	var res bool
	if *buildImagesOptsPtr == "" {
		res = RunCmdTest(t, "build", "docker", "build", "-t", imageTag, imageDir)
	} else {
		res = RunCmdTest(t, "build", "docker", "build", *buildImagesOptsPtr, "-t", imageTag, imageDir)
	}
	if !res {
		t.Fatalf("Failed to build %s image", image)
	}
	res = RunCmdTest(t, "push", "docker", "push", imageTag)
	if !res {
		t.Fatalf("Failed to push %s image", image)
	}
}

func RunCharts(t *testing.T) {
	artifacts := DiscoverArtifacts(t, "charts")
	for i, _ := range [2]struct{}{} {
		for _, chart := range artifacts {
			chartDir := path.Join(*repoPathPtr, chart)
			res := RunCmdTest(t, fmt.Sprintf("dependencies/%s/iteration_%d", chart, i+1), helmCmd, "dependency", "update", chartDir)
			if !res {
				t.Fatalf("Failed to update dependencies")
			}
		}
	}
	excludes := ListExcludesCharts()
	for _, chart := range artifacts {
		if !excludes[chart] {
			chart := chart
			t.Run(chart, func(t *testing.T) {
				t.Parallel()
				RunChart(t, chart)
			})
		}
	}
}

func RunChart(t *testing.T, chart string) {
	chartDir := path.Join(*repoPathPtr, chart)
	res := RunCmdTest(t, "lint", helmCmd, "lint", chartDir)
	if !res {
		t.Fatalf("lint failed, not proceeding")
	}
	t.Run("tests", func(t *testing.T) {
		configs, err := filepath.Glob(path.Join(*configPathPtr, chart, "*"))
		if err != nil {
			t.Fatalf("Failed to lookup configs: %s", err)
		}
		t.Logf("Found configs for chart %s: %+v", chart, configs)
		configs = append(configs, "")
		for _, config := range configs {
			var testName string
			if config != "" {
				testName = config
			} else {
				testName = "_default_"
			}
			t.Run(testName, func(t *testing.T) {
				RunOneConfig(t, chart, config)
			})
		}
	})
}

func RunPushImages(t *testing.T) {
	for _, image := range DiscoverArtifacts(t, "images") {
		image := image
		t.Run(image, func(t *testing.T) {
			t.Parallel()
			RunPushImage(t, image)
		})
	}
}

func RunPushImage(t *testing.T, image string) {
	imageDir := path.Join(*imagesPathPtr, image)
	imageVersionFilePath := path.Join(imageDir, ".version")
	versionData, err := ioutil.ReadFile(imageVersionFilePath)
	if err != nil {
		t.Fatalf("Couldn't get image version from file: %s\n%s", imageVersionFilePath, err)
	}
	version := strings.TrimSpace(string(versionData))
	imageTag := fmt.Sprintf("%s/%s:%s", *imageRepoPtr, image, version)
	newTag := fmt.Sprintf("%s/%s:%s", *pushRepoPtr, image, version)
	res := RunCmdTest(t, "tag", "docker", "tag", imageTag, newTag)
	if !res {
		t.Fatalf("Failed to tag %s image", image)
	}
	res = RunCmdTest(t, "push", "docker", "push", newTag)
	if !res {
		t.Fatalf("Failed to push %s image", image)
	}
}

func RunOneConfig(t *testing.T, chart string, config string) {
	ns := *prefixPtr + randStringRunes(10)
	rel := *prefixPtr + randStringRunes(3)
	chartDir := path.Join(*repoPathPtr, chart)
	helmHome, ok := os.LookupEnv("HELM_HOME")
	if ok {
		helmHome = helmHome + "-" + ns
	} else {
		var err error
		helmHome, err = ioutil.TempDir("", "helm-"+ns+"-")
		if err != nil {
			t.Fatalf("Failed to create temporary directory for helm home")
		}
	}

	imageParams := []string{}
	valuesData, err := ioutil.ReadFile(path.Join(chartDir, "values.yaml"))
	if err != nil {
		t.Fatalf("Failed to read values.yaml for %s chart", chart)
	}
	values := make(map[interface{}]interface{})
	err = yaml.Unmarshal(valuesData, &values)
	if err != nil {
		t.Fatalf("Error parsing values.yaml file")
	}
	image := values["image"]
	if image != nil {
		img, ok := image.(map[interface{}]interface{})
		if !ok {
			t.Fatalf("Incorrect image section\n%+v", image)
		}
		if img["repository"].(string) == "mirantisworkloads/" {
			imageParams = append(imageParams, fmt.Sprintf("image.repository=%s/", *imageRepoPtr))
		}
	}
	for key, item := range values {
		value, ok := item.(map[interface{}]interface{})
		if ok {
			image = value["image"]
			if image != nil {
				img, ok := image.(map[interface{}]interface{})
				if !ok {
					t.Fatalf("Incorrect image section\n%+v", image)
				}
				if img["repository"].(string) == "mirantisworkloads/" {
					imageParams = append(imageParams, fmt.Sprintf("%s.image.repository=%s/", key, *imageRepoPtr))
				}
			}
		}
	}

	createNsResult := RunCmdTest(t, "create_ns", kubectlCmd, "create", "ns", ns)
	if createNsResult {
		defer RunCmdTest(t, "delete_ns", kubectlCmd, "delete", "ns", ns)
	} else {
		for _, name := range []string{"install_tiller", "install", "wait_deployments", "test", "delete", "delete_ns"} {
			FailTest(t, name, "failed to create namespace")
		}
		return
	}

	installTillerResult := t.Run("install_tiller", func(t *testing.T) {
		RunCmd(t, helmCmd, "--tiller-namespace", ns, "--home", helmHome, "init")
		for i := 0; i < 10; i++ {
			cmd := exec.Command(helmCmd, "--tiller-namespace", ns, "--home", helmHome, "list")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("helm list failed: %s\nOutput: %s", err, output)
			} else {
				return
			}
			time.Sleep(3 * time.Second)
		}
		t.Fatalf("Tiller takes too long to start")
	})
	if !installTillerResult {
		for _, name := range []string{"install", "wait_deployments", "test", "delete"} {
			FailTest(t, name, "failed to init tiller")
		}
		return
	}

	installArgs := []string{helmCmd, "--tiller-namespace", ns, "--home", helmHome, "install", chartDir, "--namespace", ns, "--name", rel, "--wait", "--timeout", "600"}
	if config != "" {
		installArgs = append(installArgs, "--values", config)
	}
	if len(imageParams) > 0 {
		installArgs = append(installArgs, "--set", strings.Join(imageParams, ","))
	}
	installResult := RunCmdTest(t, "install", installArgs...)

	if !installResult {
		for _, name := range []string{"wait_deployments", "test", "delete"} {
			FailTest(t, name, "helm install failed")
		}
		return
	}

	defer RunCmdTest(t, "delete", helmCmd, "--tiller-namespace", ns, "--home", helmHome, "delete", rel, "--purge")

	// wait deployments
	waitDeploymentsResult := t.Run("wait_deployments", func(t *testing.T) {
		cmd := exec.Command(kubectlCmd, "get", "deployment", "-n", ns, "-l", fmt.Sprintf("release=%s", rel), "-o", "jsonpath={ .items[*].metadata.name }")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Get deployments command failed")
		}
		deployments := strings.Split(strings.TrimSpace(string(output)), " ")
		for _, dp := range deployments {
			if dp != "" {
				cmd := exec.Command(kubectlCmd, "rollout", "status", "-w", "-n", ns, fmt.Sprintf("deployment/%s", dp))
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("%s deployment is not ready.\nDetails:\n%s", dp, string(out))
				}
			}
		}
	})
	if !waitDeploymentsResult {
		FailTest(t, "test", "Chart is not ready")
		return
	}

	RunCmdTest(t, "test", helmCmd, "--tiller-namespace", ns, "--home", helmHome, "test", rel)
}

func CreateChartVersionVerify(t *testing.T, chartPath string, configStr string, configFile string, ns string, rel string) {
	installArgs := []string{helmCmd, "install", chartPath, "--namespace", ns, "--name", rel, "--wait", "--timeout", "600"}
	if configStr != "" {
		installArgs = append(installArgs, "--set")
		installArgs = append(installArgs, configStr)
	}
	if configFile != "" {
		installArgs = append(installArgs, "-f")
		installArgs = append(installArgs, configFile)
	}
	installResult := RunCmdTest(t, "install", installArgs...)

	if installResult {
		RunCmdTest(t, "test", helmCmd, "test", rel)
	} else {
		FailTest(t, "test", "helm install failed")
	}
}

func DeleteChartVersionVerify(t *testing.T, ns string, rel string) {
	RunCmdTest(t, "delete", helmCmd, "delete", rel, "--purge")
	RunCmdTest(t, "delete_ns", kubectlCmd, "delete", "ns", ns)
}

func RunCmdTest(t *testing.T, name string, args ...string) bool {
	return t.Run(name, func(t *testing.T) {
		RunCmd(t, args...)
	})
}

func RunCmd(t *testing.T, args ...string) {
	t.Logf("Running command: %+v", args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %s\nCommand output: %s", err, output)
	} else {
		t.Logf("Command output: %s", output)
	}
}

func FailTest(t *testing.T, name string, format string, args ...interface{}) {
	t.Run(name, func(t *testing.T) { t.Fatalf(format, args...) })
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
