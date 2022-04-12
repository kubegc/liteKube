package app

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
	"github.com/litekube/LiteKube/pkg/options/leader"
	options "github.com/litekube/LiteKube/pkg/options/leader"
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
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested() // --version=false/true/simple/raw to print version

			klog.Info("Welcome to LiteKube Leader, a lite cluster with almost nearly the same capabilities as the K8S Leader node")
			klog.Info(version.GetSimple())

			// // load config from disk-file and merge with flags
			// if errs := opt.LoadConfig(); len(errs) != 0 {
			// 	klog.Error("some error in your configs")
			// 	return fmt.Errorf("some error in your configs")
			// }

			// // complete all default server options,current is none-function
			// if err := opt.Complete(); err != nil {
			// 	klog.Errorf("complete options error: %v", err)
			// 	return err
			// }

			// // ready to run
			// return Run(opt, util.SetupSignalHandler())
			return nil
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
	yamlFmt := "\n[config-file format]:"
	flagSections := opt.HelpSections()
	helpSection := opt.ConfigHelpSection()
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		for _, section := range flagSections {
			section.PrintSection(cmd.OutOrStderr(), help.FormatClamp("<", ">"))
		}
		fmt.Fprintln(cmd.OutOrStderr(), yamlFmt)
		helpSection.PrintSection(cmd.OutOrStderr(), help.FormatHeader("// "))
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		for _, section := range flagSections {
			section.PrintSection(cmd.OutOrStderr(), help.FormatClamp("<", ">"))
		}
		fmt.Fprintln(cmd.OutOrStderr(), yamlFmt)
		helpSection.PrintSection(cmd.OutOrStderr(), help.FormatHeader("// "))
	})

	return cmd

}

func addFlags(cmd *cobra.Command) {
	leader.AddFlagsTo(cmd.Flags())
	verflag.AddFlagsTo(cmd.Flags())
}
