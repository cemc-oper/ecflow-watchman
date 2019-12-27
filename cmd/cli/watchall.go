package cli

import (
	"github.com/nwpc-oper/ecflow-watchman"
	"github.com/pquerna/ffjson/ffjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	configFilePath   = ""
	isProfiling      = true
	profilingAddress = "127.0.0.1:30485"
)

func init() {
	rootCmd.AddCommand(watchAllCmd)

	watchAllCmd.Flags().StringVar(&configFilePath, "config-file", "", "config file path, required")
	watchAllCmd.Flags().BoolVar(&isProfiling, "enable-profiling", false, "enable profiling")
	watchAllCmd.Flags().StringVar(&profilingAddress, "profiling-address", "127.0.0.1:30485", "profiling address")
	watchAllCmd.MarkFlagRequired("config-file")
}

const watchAllCommandDescription = `
watch all ecFlow servers listed in the configure file.

For each ecflow server, watch-all will store status in redis with some key and publish status with some channel.

Config file is as follows:

global:
  scrape_interval: 20s
  scrape_timeout: 10s # not worked

scrape_configs:
  -
    job_name: job name
    owner: owner
    repo: repo
    host: ecflow server host
    port: ecflow server port

sink_config:
  type: redis # only redis is supported
  url: redis url
`

var watchAllCmd = &cobra.Command{
	Use:   "watch-all",
	Short: "watch all ecFlow servers",
	Long:  watchAllCommandDescription,
	Run: func(cmd *cobra.Command, args []string) {
		if isProfiling {
			log.Infof("enable profiling...")
			go func() {
				log.Println(http.ListenAndServe(profilingAddress, nil))
			}()
		}

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

		// create redis client
		storer := ecflow_watchman.RedisStorer{
			Address:  redisUrl,
			Password: "",
			Database: 0,
		}
		storer.Create()
		defer storer.Close()

		// create publisher
		redisPublisher := ecflow_watchman.RedisPublisher{
			Client:   nil,
			Address:  redisUrl,
			Password: "",
			Database: 1,
		}
		redisPublisher.Create()
		defer redisPublisher.Close()

		for _, job := range config.ScrapeConfigs {
			// create redis publisher for a scrape job
			//messages := make(chan []byte)
			//go redisPub(job, &redisPublisher, messages)

			// create collect goroutine for a scrape job
			go scrapeStatus(job, &storer, scrapeInterval)

			log.WithFields(log.Fields{
				"owner": job.Owner,
				"repo":  job.Repo,
			}).Info("new job loaded")
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

func scrapeStatus(job ScrapeJob, storer ecflow_watchman.Storer, scrapeInterval time.Duration) {
	c := time.Tick(scrapeInterval)
	for _ = range c {
		// get ecflow server status
		ecflowServerStatus := ecflow_watchman.GetEcflowStatus(job.EcflowServerConfig)
		if ecflowServerStatus == nil {
			// ignore any error,
			// continue to next loop when we can't get ecflow status.
			continue
		}

		b, err := ffjson.Marshal(ecflowServerStatus)
		//ecflowServerStatus = nil
		if err != nil {
			log.WithFields(log.Fields{
				"owner": job.EcflowServerConfig.Owner,
				"repo":  job.EcflowServerConfig.Repo,
			}).Errorf("Marshal json has error: %v", err)
			continue
		}

		// save to redis
		storer.Send(job.EcflowServerConfig.Owner, job.EcflowServerConfig.Repo, b)

		// send message to channel
		//messages <- b

		//ffjson.Pool(b)
		b = nil
	}
}

func redisPub(job ScrapeJob, redisPublisher *ecflow_watchman.RedisPublisher, messages chan []byte) {
	channelName := job.Owner + "/" + job.Repo + "/status/channel"

	redisPublisher.CreatePubsub(channelName)
	log.WithFields(log.Fields{
		"owner": job.Owner,
		"repo":  job.Repo,
	}).Infof("subscribe redis...%s", channelName)

	for message := range messages {
		log.WithFields(log.Fields{
			"owner": job.Owner,
			"repo":  job.Repo,
		}).Infof("publish to redis...")

		err := redisPublisher.Publish(channelName, message)
		message = nil

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
}
