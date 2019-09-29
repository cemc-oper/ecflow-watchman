package cli

import (
	ecflow_watchman "ecflow-watchman"
	"github.com/spf13/cobra"
	"time"
)

var (
	ecflowHost = ""
	ecflowPort = ""
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVar(&ecflowHost, "ecflow-host", "", "ecFlow server host")
	watchCmd.Flags().StringVar(&ecflowPort, "ecflow-port", "", "ecFlow server port")
	watchCmd.MarkFlagRequired("ecflow-host")
	watchCmd.MarkFlagRequired("ecflow-port")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch ecFlow servers",
	Long:  "watch ecFlow servers",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			ecflow_watchman.GetEcflowStatus(ecflowHost, ecflowPort)
			time.Sleep(30 * time.Second)
		}
	},
}
