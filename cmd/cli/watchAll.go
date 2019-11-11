package cli

import (
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
			go func(job ScrapeJob, redisUrl string) {

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

					// save to redis key
					go func(
						config ecflow_watchman.EcflowServerConfig,
						ecflowServerStatus ecflow_watchman.EcflowServerStatus,
						redisUrl string) {
						ecflow_watchman.StoreToRedis(config, ecflowServerStatus, redisUrl)
					}(job.EcflowServerConfig, *ecflowServerStatus, redisUrl)
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
