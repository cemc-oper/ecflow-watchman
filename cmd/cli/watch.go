package cli

import (
	"encoding/json"
	"github.com/perillaroc/ecflow-watchman"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var (
	owner          = ""
	repo           = ""
	ecflowHost     = ""
	ecflowPort     = ""
	redisUrl       = ""
	scrapeInterval = ""
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVar(&owner, "owner", "", "owner")
	watchCmd.Flags().StringVar(&repo, "repo", "", "repo")
	watchCmd.Flags().StringVar(&ecflowHost, "ecflow-host", "", "ecFlow server host")
	watchCmd.Flags().StringVar(&ecflowPort, "ecflow-port", "", "ecFlow server port")
	watchCmd.Flags().StringVar(&redisUrl, "redis-url", "", "redis url")
	watchCmd.Flags().StringVar(&scrapeInterval, "scrape-interval", "10s", "scrape interval")
	watchCmd.MarkFlagRequired("owner")
	watchCmd.MarkFlagRequired("port")
	watchCmd.MarkFlagRequired("ecflow-host")
	watchCmd.MarkFlagRequired("ecflow-port")
	watchCmd.MarkFlagRequired("redis-url")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch a single ecFlow server",
	Long:  "watch a single ecFlow server",
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := time.ParseDuration(scrapeInterval)
		if err != nil {
			panic(err)
		}
		c := time.Tick(duration)
		for _ = range c {
			config := ecflow_watchman.EcflowServerConfig{
				Owner:          owner,
				Repo:           repo,
				Host:           ecflowHost,
				Port:           ecflowPort,
				ConnectTimeout: 10,
			}
			ecflowServerStatus := ecflow_watchman.GetEcflowStatus(config)
			if ecflowServerStatus == nil {
				continue
			}
			b, err := json.Marshal(ecflowServerStatus)
			if err != nil {
				log.WithFields(log.Fields{
					"owner": config.Owner,
					"repo":  config.Repo,
				}).Error("Marshal json has error: ", err)
				return
			}
			ecflow_watchman.StoreToRedis(config, b, redisUrl)
		}
	},
}
