package ecflow_watchman

import (
	"encoding/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

func StoreToRedis(config EcflowServerConfig, ecflowServerStatus EcflowServerStatus, redisUrl string) {
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
