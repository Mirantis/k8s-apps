package main

import (
	"fmt"
	"os"

	"github.com/Mirantis/k8s-apps/helm-apply/apply"
	"github.com/Mirantis/k8s-apps/helm-apply/getter"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
)

var (
	verbose bool
	diff    bool
	dryRun  bool
)

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	applier, err := apply.NewApplier(args[0], getter.Providers())
	if err != nil {
		return err
	}
	err = applier.Run(verbose, diff, dryRun)
	if err != nil {
		return err
	}
	return nil
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("configuration file path is required")
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:     "apply CONFIG",
		Short:   "Apply declaratively defined configuration of helm charts",
		PreRunE: validateFlags,
		RunE:    run,
	}
	runtime.ErrorHandlers = runtime.ErrorHandlers[1:]
	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "show release description after deployment")
	f.BoolVar(&diff, "diff", false, "show diff between deployed and applied configuration")
	f.BoolVar(&dryRun, "dry-run", false, "process charts without actual deployment")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
