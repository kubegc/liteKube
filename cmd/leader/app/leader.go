package app

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/leader/config"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	"github.com/litekube/LiteKube/pkg/util"
	"github.com/litekube/LiteKube/pkg/version"
	verflag "github.com/litekube/LiteKube/pkg/version/varflag"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var ComponentName = "Leader"

func NewLeaderCommand() *cobra.Command {
	opt := options.NewLeaderOptions()

	cmd := &cobra.Command{
		Use:  ComponentName,
		Long: fmt.Sprintf("%s is a lite leader-component with almost nearly the same capabilities as the K8S Leader", ComponentName),

		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested() // --version=false/true/simple/raw to print version

			klog.Infof("Welcome to LiteKube %s, a lite cluster with almost nearly the same capabilities as the K8S Leader node", ComponentName)
			klog.Info(version.GetSimple())

			// load config from --config-file
			if err := opt.LoadConfig(); err != nil {
				return err
			}

			// run leader
			return Run(opt, util.SetupSignalHandler())
		},
		Args: func(cmd *cobra.Command, args []string) error { // Validate unresolved args
			for _, arg := range args {
				if len(arg) > 0 {
					klog.Errorf("%q does not support subcommands at this time but get %q", cmd.CommandPath(), args)
					return fmt.Errorf("%q does not support subcommands at this time but get %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// add flags to cmd
	addFlags(cmd)

	// add help tips for program
	usageFmt := "Usage:\n  %s\n\n"
	yamlFmt := "\n[config-file format]:\n// setting for kube-apiserver,kube-controller-manager,kube-scheduler and litekube additions\n"
	flagSections := opt.HelpSections()
	yamlSection := opt.ConfigHelpSection()
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		// print flags help
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		for _, section := range flagSections {
			section.PrintSection(cmd.OutOrStderr(), help.FormatClamp("<", ">"))
		}

		// print yaml help
		fmt.Fprintln(cmd.OutOrStderr(), yamlFmt)
		for _, section := range yamlSection {
			section.PrintSection(cmd.OutOrStderr(), help.FormatHeader("# "))
			fmt.Fprintln(cmd.OutOrStderr())
		}
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// print flags help
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		for _, section := range flagSections {
			section.PrintSection(cmd.OutOrStderr(), help.FormatClamp("<", ">"))
		}

		// print yaml help
		fmt.Fprintln(cmd.OutOrStderr(), yamlFmt)
		for _, section := range yamlSection {
			section.PrintSection(cmd.OutOrStderr(), help.FormatHeader("# "))
			fmt.Fprintln(cmd.OutOrStderr())
		}
	})

	return cmd

}

func Run(opt *options.LeaderOptions, stopCh <-chan struct{}) error {
	runtimeConfig := config.NewLeaderRuntime(opt)
	defer runtimeConfig.Stop()

	if err := runtimeConfig.LoadFlags(); err != nil {
		return err
	}

	// print finally runtime-flags
	runtimeConfig.RuntimeOption.PrintFlags(func() func(format string, a ...interface{}) error {
		return func(format string, a ...interface{}) error {
			klog.Infof(format, a...)
			return nil
		}
	}())

	// run k8s
	if err := runtimeConfig.Run(); err != nil {
		return err
	}

	<-stopCh // wait util read system close signal

	klog.Info("We have prepare to close process, it won't take you too much time, wait please!")

	return nil
}

func addFlags(cmd *cobra.Command) {
	options.AddFlagsTo(cmd.Flags())
	//verflag.AddFlagsTo(cmd.Flags())
}
