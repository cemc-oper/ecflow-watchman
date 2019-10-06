package ecflow_watchman

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/perillaroc/ecflow-client-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type EcflowServerStatus struct {
	StatusRecords []ecflow_client.StatusRecord `json:"status_records"`
	CollectedTime time.Time                    `json:"collected_time"`
}

type EcflowServerConfig struct {
	Owner          string `yaml:"owner"`
	Repo           string `yaml:"repo"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	ConnectTimeout int    `yaml:"connect_timeout"`
}

func GetEcflowStatus(config EcflowServerConfig, redisUrl string) {
	client := ecflow_client.CreateEcflowClient(config.Host, config.Port)
	client.SetConnectTimeout(config.ConnectTimeout)
	defer client.Close()

	ret := client.Sync()
	if ret != 0 {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Error("sync has error: ", ret)
		return
	}

	records := client.StatusRecords()

	ecflowServerStatus := EcflowServerStatus{
		StatusRecords: records,
		CollectedTime: client.CollectedTime,
	}

	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Info(
		"get ",
		len(ecflowServerStatus.StatusRecords),
		" nodes at ",
		ecflowServerStatus.CollectedTime.Format("2006-01-02 15:04:05.999999"))

	b, err := json.Marshal(ecflowServerStatus)
	if err != nil {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Error("Marshal json has error: ", err)
		return
	}

	key := config.Owner + "/" + config.Repo + "/status"

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	err = redisClient.Set(key, b, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Error("store to redis has error: ", err)
		return
	}

	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Info(
		"write to redis at ",
		time.Now().Format("2006-01-02 15:04:05.999999"))

}
