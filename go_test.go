package go_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	paramsPtr          = flag.String("params", "", "Set config values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	excludePtr         = flag.String("exclude", "", "List of charts to exclude from run")
	prefixPtr          = flag.String("prefix", "", "Prefix to prepend to object names (releases, namespaces)")
	chartsPtr          = flag.Bool("charts", true, "Test charts")
	imagesPtr          = flag.Bool("images", false, "Build images")
	pushPtr            = flag.Bool("push", false, "Push images after tests")
	imageRepoPtr       = flag.String("image-repo", "127.0.0.1:5000", "Image repo address for test")
	pushRepoPtr        = flag.String("push-repo", "mirantisworkloads", "Image repo address for push")
	buildImagesOptsPtr = flag.String("build-images-opts", "", "Docker opts for building images")
	helmCmd            = LookupEnvDefault("HELM_CMD", "helm")
	kubectlCmd         = LookupEnvDefault("KUBECTL_CMD", "kubectl")
	excludes           []string
)

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
	excludes := make(map[string]bool)
	for _, e := range strings.Split(*excludePtr, ",") {
		excludes[filepath.Base(e)] = true
	}
	var res []string
	for _, a := range artifacts {
		if !excludes[a] {
			res = append(res, a)
		}
	}
	return res
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
	for _, chart := range DiscoverArtifacts(t, "charts") {
		chart := chart
		t.Run(chart, func(t *testing.T) {
			t.Parallel()
			RunChart(t, chart)
		})
	}
}

func RunChart(t *testing.T, chart string) {
	chartDir := path.Join(*repoPathPtr, chart)
	res := RunCmdTest(t, "lint", helmCmd, "lint", chartDir)
	if !res {
		t.Fatalf("lint failed, not proceeding")
	}
	res = RunCmdTest(t, "dependencies", helmCmd, "dependency", "update", chartDir)
	if !res {
		t.Fatalf("Failed to update dependencies")
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
			RunImage(t, image)
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
		for _, name := range []string{"install_tiller", "install", "test", "delete", "delete_ns"} {
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
		for _, name := range []string{"install", "test", "delete"} {
			FailTest(t, name, "failed to init tiller")
		}
		return
	}

	installArgs := []string{helmCmd, "--tiller-namespace", ns, "--home", helmHome, "install", chartDir, "--namespace", ns, "--name", rel, "--wait", "--timeout", "600"}
	if config != "" {
		installArgs = append(installArgs, "--values", config)
	}
	if *paramsPtr != "" {
		installArgs = append(installArgs, "--set", *paramsPtr)
	}
	installResult := RunCmdTest(t, "install", installArgs...)

	if installResult {
		RunCmdTest(t, "test", helmCmd, "--tiller-namespace", ns, "--home", helmHome, "test", rel)
	} else {
		FailTest(t, "test", "helm install failed")
	}

	RunCmdTest(t, "delete", helmCmd, "--tiller-namespace", ns, "--home", helmHome, "delete", rel, "--purge")
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
