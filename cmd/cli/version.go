package cli

import (
	"github.com/perillaroc/ecflow-watchman"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		ecflow_watchman.PrintVersionInformation()
	},
}
