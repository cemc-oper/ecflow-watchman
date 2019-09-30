package cli

import (
	"github.com/perillaroc/ecflow-watchman"
	"github.com/spf13/cobra"
	"time"
)

var (
	owner      = ""
	repo       = ""
	ecflowHost = ""
	ecflowPort = ""
	redisUrl   = ""
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVar(&owner, "owner", "", "owner")
	watchCmd.Flags().StringVar(&repo, "repo", "", "repo")
	watchCmd.Flags().StringVar(&ecflowHost, "ecflow-host", "", "ecFlow server host")
	watchCmd.Flags().StringVar(&ecflowPort, "ecflow-port", "", "ecFlow server port")
	watchCmd.Flags().StringVar(&redisUrl, "redis-url", "", "redis url")
	watchCmd.MarkFlagRequired("owner")
	watchCmd.MarkFlagRequired("port")
	watchCmd.MarkFlagRequired("ecflow-host")
	watchCmd.MarkFlagRequired("ecflow-port")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch ecFlow servers",
	Long:  "watch ecFlow servers",
	Run: func(cmd *cobra.Command, args []string) {
		c := time.Tick(10 * time.Second)
		for _ = range c {
			ecflow_watchman.GetEcflowStatus(owner, repo, ecflowHost, ecflowPort, redisUrl)
		}
	},
}
