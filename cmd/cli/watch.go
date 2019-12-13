package cli

import (
	"encoding/json"
	"github.com/nwpc-oper/ecflow-watchman"
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

	watchCmd.Flags().StringVar(&owner, "owner", "", "owner name for redis key")
	watchCmd.Flags().StringVar(&repo, "repo", "", "repo name for redis key")
	watchCmd.Flags().StringVar(&ecflowHost, "ecflow-host", "", "ecFlow server host")
	watchCmd.Flags().StringVar(&ecflowPort, "ecflow-port", "", "ecFlow server port")
	watchCmd.Flags().StringVar(&redisUrl, "redis-url", "", "redis url")
	watchCmd.Flags().StringVar(&scrapeInterval, "scrape-interval", "10s",
		"scrape interval, time duration such as 5s.")

	watchCmd.MarkFlagRequired("owner")
	watchCmd.MarkFlagRequired("port")
	watchCmd.MarkFlagRequired("ecflow-host")
	watchCmd.MarkFlagRequired("ecflow-port")
	watchCmd.MarkFlagRequired("redis-url")
}

const watchCommandDescription = `
Watch a single ecFlow server and save its status into redis.
`

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch a single ecFlow server",
	Long:  watchCommandDescription,
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := time.ParseDuration(scrapeInterval)
		if err != nil {
			panic(err)
		}

		config := ecflow_watchman.EcflowServerConfig{
			Owner:          owner,
			Repo:           repo,
			Host:           ecflowHost,
			Port:           ecflowPort,
			ConnectTimeout: 10,
		}

		// create redis client
		storer := ecflow_watchman.RedisStorer{
			Address:  redisUrl,
			Password: "",
			Database: 0,
		}
		storer.Create()
		defer storer.Close()

		c := time.Tick(duration)
		for _ = range c {
			ecflowServerStatus := ecflow_watchman.GetEcflowStatus(config)
			if ecflowServerStatus == nil {
				log.WithFields(log.Fields{
					"owner":     config.Owner,
					"repo":      config.Repo,
					"component": "watch",
				}).Error("get ecflow status has error.")
				continue
			}

			b, err := json.Marshal(ecflowServerStatus)
			if err != nil {
				log.WithFields(log.Fields{
					"owner":     config.Owner,
					"repo":      config.Repo,
					"component": "watch",
				}).Errorf("encode json has error: %v", err)
				return
			}

			storer.Send(config.Owner, config.Repo, b)
		}
	},
}
