package go_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/kubernetes/apimachinery/pkg/util/yaml"
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
	repoPathPtr         = flag.String("repo", "charts/", "Path to charts repository")
	imagesPathPtr       = flag.String("images-dir", "images/", "Path to Dockerfiles")
	configPathPtr       = flag.String("config", "tests/", "Path to charts config files")
	excludePtr          = flag.String("exclude", "", "List of charts to exclude from run")
	prefixPtr           = flag.String("prefix", "", "Prefix to prepend to object names (releases, namespaces)")
	chartsPtr           = flag.Bool("charts", true, "Test charts")
	kubernetesDomainPtr = flag.String("kubernetes-domain", "cluster.local", "Base domain of Kubernetes cluster tests are run on")
	imagesPtr           = flag.Bool("images", false, "Build images")
	pushPtr             = flag.Bool("push", false, "Push images after tests")
	imageRepoPtr        = flag.String("image-repo", "127.0.0.1:5000", "Image repo address for test")
	pushRepoPtr         = flag.String("push-repo", "mirantisworkloads", "Image repo address for push")
	buildImagesOptsPtr  = flag.String("build-images-opts", "", "Docker opts for building images")
	verifyVersion       = flag.Bool("verify-version", false, "Run tests to verify new helm/k8s version")
	remoteCluster       = flag.String("remote-cluster", "127.0.0.1", "Cluster IP if remote kubernetes cluster is used")
	verifyIngress       = flag.Bool("verify-ingress", false, "Ensure ingress is working correctly")
	ingressSvc          = flag.String("ingress-svc", "", "Ingress service host:port to connect the ingress resource")
	helmCmd             = LookupEnvDefault("HELM_CMD", "helm")
	kubectlCmd          = LookupEnvDefault("KUBECTL_CMD", "kubectl")
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

func DiscoverArtifacts(t *testing.T, dir string) []string {
	matches, err := filepath.Glob(dir + "/*")
	if err != nil {
		t.Fatalf("Failed to list directory %s: %s", dir, err)
	}
	var allArtifacts []string
	for _, match := range matches {
		_, err := os.Stat(path.Join(match, ".nobuild"))
		if !os.IsNotExist(err) {
			continue
		}
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
		configs, err := filepath.Glob(path.Join(*configPathPtr, chart, "*.yaml"))
		if err != nil {
			t.Fatalf("Failed to lookup configs: %s", err)
		}
		if len(configs) == 0 {
			t.Fatalf("Didn't find any configs to run tests for chart '%s'", chart)
		}
		t.Logf("Found configs for chart %s: %+v", chart, configs)
		for _, config := range configs {
			testName := path.Base(config)
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
		for i := 0; i < 20; i++ {
			cmd := exec.Command(helmCmd, "--tiller-namespace", ns, "--home", helmHome, "list")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("helm list failed: %s\nOutput: %s", err, output)
			} else {
				return
			}
			time.Sleep(6 * time.Second)
		}
		t.Fatalf("Tiller takes too long to start")
	})
	if !installTillerResult {
		for _, name := range []string{"install", "wait_deployments", "test", "delete"} {
			FailTest(t, name, "failed to init tiller")
		}
		return
	}

	installResult := t.Run("install", func(t *testing.T) {
		RunHelmInstall(t, helmHome, ns, rel, chartDir, config)
	})

	if !installResult {
		for _, name := range []string{"wait_deployments", "test", "delete"} {
			FailTest(t, name, "helm install failed")
		}
		return
	}

	defer RunCmdTest(t, "delete", helmCmd, "--tiller-namespace", ns, "--home", helmHome, "delete", rel, "--purge")

	// wait deployments
	waitDeploymentsResult := t.Run("wait_deployments", func(t *testing.T) {
		WaitForResources(t, helmHome, ns, rel)
	})
	if !waitDeploymentsResult {
		FailTest(t, "test", "Chart is not ready")
		return
	}

	t.Run("test", func(t *testing.T) {
		RunHelmTest(t, helmHome, ns, rel)
	})
}

func RunHelmInstall(t *testing.T, helmHome string, ns string, rel string, chartDir string, config string) {
	installArgs := []string{helmCmd, "--tiller-namespace", ns, "--home", helmHome, "install", chartDir, "--namespace", ns, "--name", rel, "--wait", "--timeout", "600"}
	installArgs = append(installArgs, "--values", "/dev/stdin")
	tmpl, err := template.ParseFiles(config)
	if err != nil {
		t.Fatalf("Failed to parse template %s: %s", config, err)
	}
	var buf bytes.Buffer
	obj := struct {
		Repository       string
		KubernetesDomain string
	}{*imageRepoPtr + "/", *kubernetesDomainPtr}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		t.Fatalf("Failed to execute template %s: %s", config, err)
	}
	str_config := buf.String()
	cmd := exec.Command(installArgs[0], installArgs[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to get stding pipe: %s", err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, str_config)
	}()
	t.Logf("Running command: %+v", installArgs)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %s\nCommand output: %s", err, output)
	}
	t.Logf("Command output: %s", output)
}

func WaitForResources(t *testing.T, helmHome string, ns string, rel string) {
	cmd := exec.Command(helmCmd, "--tiller-namespace", ns, "--home", helmHome, "get", "manifest", rel)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get manifest: %s", err)
	}
	deployments := make(map[string]bool)
	var obj struct {
		Kind     string `json:"kind"`
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(output)), 1)
	for {
		err := decoder.Decode(&obj)
		if err != nil {
			if err != io.EOF {
				t.Fatalf("Failed to parse manifest: %s", err)
			}
			break
		}
		if obj.Kind == "Deployment" {
			deployments[obj.Metadata.Name] = true
		}
	}
	t.Logf("Found deployments: %l", deployments)
	if len(deployments) == 0 {
		return
	}
	cmd = exec.Command(kubectlCmd, "get", "deployment", "-n", ns, "-o", "json", "-w")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe from kubectl to watch deployments: %s", err)
	}
	decoder = yaml.NewYAMLOrJSONDecoder(stdout, 1)
	var deployment struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Status struct {
			Replicas      int `json:"replicas"`
			ReadyReplicas int `json:"readyReplicas"`
		} `json:"status"`
	}
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start kubectl to watch deployments: %s", err)
	}
	for {
		err := decoder.Decode(&deployment)
		if err != nil {
			t.Logf("Failed to decode JSON from kubectl: %s", err)
			t.Fail()
			break
		}
		t.Logf("%s replicas=%d ready=%d\n", deployment.Metadata.Name, deployment.Status.Replicas, deployment.Status.ReadyReplicas)
		if deployment.Status.Replicas != 0 && deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			delete(deployments, deployment.Metadata.Name)
			if len(deployments) == 0 {
				break
			}
		}
	}
	err = cmd.Process.Kill()
	if err != nil {
		t.Fatalf("Failed to kill kubectl process: %s", err)
	}
}

func RunHelmTest(t *testing.T, helmHome string, ns string, rel string) {
	args := []string{helmCmd, "--tiller-namespace", ns, "--home", helmHome, "test", rel}
	t.Logf("Running command: %+v", args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command failed: %s\nCommand output: %s", err, output)
		re := regexp.MustCompile("`(kubectl logs.*)`")
		for _, match := range re.FindAllSubmatch(output, -1) {
			match_str := string(match[1])
			args := strings.Split(match_str, " ")
			if args[0] == "kubectl" {
				args[0] = kubectlCmd
			}
			cmd := exec.Command(args[0], args[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Failed to get logs from `%s`: %s\nCommand output: %s", match_str, err, output)
			} else {
				t.Logf("Output from `%s`:\n%s", match_str, output)
			}
		}
		t.Fail()
	} else {
		t.Logf("Command output: %s", output)
	}
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
