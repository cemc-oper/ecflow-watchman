package cli

import (
	"encoding/json"
	"github.com/perillaroc/ecflow-watchman"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	configFilePath = ""
)

func init() {
	rootCmd.AddCommand(watchAllCmd)

	watchAllCmd.Flags().StringVar(&configFilePath, "config-file", "", "config file path")
	watchAllCmd.MarkFlagRequired("config-file")
}

var watchAllCmd = &cobra.Command{
	Use:   "watch-all",
	Short: "watch all ecFlow servers",
	Long:  "watch all ecFlow servers",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readConfig(configFilePath)
		if err != nil {
			panic(err)
		}

		// scrape interval
		scrapeInterval, err := time.ParseDuration(config.Global.ScrapeInterval)
		if err != nil {
			panic(err)
		}
		log.Info("scrape_interval = ", scrapeInterval)

		// sink
		sinkType := config.SinkConfig["type"].(string)
		if sinkType != "redis" {
			log.Fatal("sink type is not supported: ", sinkType)
		}
		redisUrl := config.SinkConfig["url"].(string)
		log.Info("sink to redis: ", redisUrl)

		for _, job := range config.ScrapeConfigs {
			// create redis publisher for each scrape job
			messages := make(chan []byte)
			go func(job ScrapeJob, redisUrl string) {
				channelName := job.Owner + "/" + job.Repo + "/status/channel"

				redisPublisher := ecflow_watchman.RedisPublisher{
					Client:      nil,
					Pubsub:      nil,
					ChannelName: channelName,
					Address:     redisUrl,
					Password:    "",
					Database:    0,
				}

				redisPublisher.Create()
				log.WithFields(log.Fields{
					"owner": job.Owner,
					"repo":  job.Repo,
				}).Infof("subscribe redis...%s", channelName)
				defer redisPublisher.Close()

				for message := range messages {
					log.WithFields(log.Fields{
						"owner": job.Owner,
						"repo":  job.Repo,
					}).Infof("publish to redis...")

					err = redisPublisher.Publish(message)

					if err != nil {
						log.WithFields(log.Fields{
							"owner": job.Owner,
							"repo":  job.Repo,
						}).Errorf("publish to redis has error: %v", err)

					} else {
						log.WithFields(log.Fields{
							"owner": job.Owner,
							"repo":  job.Repo,
						}).Infof("publish to redis...done")
					}
				}
			}(job, redisUrl)

			// create collect goroutine for each scrape job
			go func(job ScrapeJob, redisUrl string, scrapeInterval time.Duration) {
				c := time.Tick(scrapeInterval)
				for _ = range c {
					// get ecflow server status
					ecflowServerStatus := ecflow_watchman.GetEcflowStatus(job.EcflowServerConfig)
					if ecflowServerStatus == nil {
						return
					}

					b, err := json.Marshal(ecflowServerStatus)
					if err != nil {
						log.WithFields(log.Fields{
							"owner": job.EcflowServerConfig.Owner,
							"repo":  job.EcflowServerConfig.Repo,
						}).Error("Marshal json has error: ", err)
						return
					}

					// save to redis key
					ecflow_watchman.StoreToRedis(job.EcflowServerConfig, b, redisUrl)

					// NOTE: may cause I/O timeout when running in goroutine.
					//go func(
					//	config ecflow_watchman.EcflowServerConfig,
					//	message []byte,
					//	redisUrl string) {
					//	ecflow_watchman.StoreToRedis(config, message, redisUrl)
					//}(job.EcflowServerConfig, *ecflowServerStatus, redisUrl)

					// send message to channel
					messages <- b
				}
			}(job, redisUrl, scrapeInterval)

			log.Info("new job loaded: ", job.Owner, "/", job.Repo)
		}

		// block forever in the main goroutine
		// see: https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
		select {}
	},
}

type ConfigDict map[interface{}]interface{}
type ConfigArray []interface{}

func readConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type Config struct {
	Global        GlobalConfig `yaml:"global"`
	ScrapeConfigs []ScrapeJob  `yaml:"scrape_configs"`
	SinkConfig    ConfigDict   `yaml:"sink_config"`
}

type ScrapeJob struct {
	JobName                            string `yaml:"job_name"`
	ecflow_watchman.EcflowServerConfig `yaml:",inline"`
}

type GlobalConfig struct {
	ScrapeInterval string `yaml:"scrape_interval"`
	ScrapeTimeout  string `yaml:"scrape_timeout"`
}
