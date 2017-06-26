package go_test

import (
	"flag"
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
	repoPathPtr   = flag.String("repo", "charts/", "Path to charts repository")
	junitPathPtr  = flag.String("junit", "report.xml", "Path to output junit-xml report")
	configPathPtr = flag.String("config", "tests/", "Path to charts config files")
	paramsPtr     = flag.String("params", "", "Set config values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	excludePtr    = flag.String("exclude", "", "List of charts to exclude from run")
	prefixPtr     = flag.String("prefix", "", "Prefix to prepend to object names (releases, namespaces)")
	helmCmd       = LookupEnvDefault("HELM_CMD", "helm")
	kubectlCmd    = LookupEnvDefault("KUBECTL_CMD", "kubectl")
	excludes      []string
)

func TestRoot(t *testing.T) {
	t.Run("charts", RunCharts)
}

func RunCharts(t *testing.T) {
	matches, err := filepath.Glob(*repoPathPtr + "/*")
	if err != nil {
		t.Fatalf("Failed to list directory %s: %s", *repoPathPtr, err)
	}
	t.Logf("Found charts: %+v", matches)
	var charts []string
	if flag.NArg() == 0 {
		for _, match := range matches {
			charts = append(charts, filepath.Base(match))
		}
	} else {
		matchSet := make(map[string]bool)
		for _, match := range matches {
			matchSet[filepath.Base(match)] = true
		}
		var not_found []string
		args := flag.Args()
		for _, arg := range args {
			if !matchSet[arg] {
				not_found = append(not_found, arg)
			}
		}
		if len(not_found) > 0 {
			t.Fatalf("Couldn't find these charts: %+v", not_found)
		}
		charts = args
	}

	excludes := make(map[string]bool)
	for _, e := range strings.Split(*excludePtr, ",") {
		excludes[e] = true
	}

	for _, chart := range charts {
		if excludes[chart] {
			continue
		}
		t.Run(chart, func(t *testing.T) { RunChart(t, chart) })
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

func RunOneConfig(t *testing.T, chart string, config string) {
	ns := *prefixPtr + randStringRunes(10)
	rel := *prefixPtr + randStringRunes(3)
	chartDir := path.Join(*repoPathPtr, chart)

	installArgs := []string{helmCmd, "install", chartDir, "--namespace", ns, "--name", rel, "--wait", "--timeout", "600"}
	if config != "" {
		installArgs = append(installArgs, "--values", config)
	}
	if *paramsPtr != "" {
		installArgs = append(installArgs, "--set", *paramsPtr)
	}
	installResult := RunCmdTest(t, "install", installArgs...)

	if installResult {
		RunCmdTest(t, "test", helmCmd, "test", rel)
	} else {
		FailTest(t, "test", "helm install failed")
	}

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
