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

func GetEcflowStatus(owner string, repo string, host string, port string, redisUrl string) {
	client := ecflow_client.CreateEcflowClient(host, port)
	defer client.Close()

	ret := client.Sync()
	if ret != 0 {
		log.Fatal("sync has error")
	}

	records := client.StatusRecords()

	ecflowServerStatus := EcflowServerStatus{
		StatusRecords: records,
		CollectedTime: client.CollectedTime,
	}

	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info(
		"get ",
		len(ecflowServerStatus.StatusRecords),
		" nodes at ",
		ecflowServerStatus.CollectedTime.Format("2006-01-02 15:04:05.999999"))

	b, err := json.Marshal(ecflowServerStatus)
	if err != nil {
		log.WithFields(log.Fields{
			"owner": owner,
			"repo":  repo,
		}).Error("Marshal json has error: ", err)
		return
	}

	key := owner + "/" + repo + "/status"

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	err = redisClient.Set(key, b, 0).Err()
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info(
		"write to redis at ",
		ecflowServerStatus.CollectedTime.Format("2006-01-02 15:04:05.999999"))

}
