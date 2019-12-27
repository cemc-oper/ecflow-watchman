package ecflow_watchman

import (
	"bytes"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Storer interface {
	Create()
	Send(owner string, repo string, message *bytes.Buffer)
	Close()
}

type RedisStorer struct {
	Address  string
	Password string
	Database int
	client   *redis.Client
}

func (s *RedisStorer) Create() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.Address,
		Password: s.Password,
		DB:       s.Database,
	})
}

func (s *RedisStorer) Close() {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}
}

func (s *RedisStorer) Send(owner string, repo string, message *bytes.Buffer) {
	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Infof("store to redis... ")

	key := owner + "/" + repo + "/status"
	err := s.client.Set(key, message.String(), 0).Err()

	if err != nil {
		log.WithFields(log.Fields{
			"owner": owner,
			"repo":  repo,
		}).Error("store to redis has error: ", err)
		return
	}

	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info("store to redis...done")
}

func StoreToRedis(config EcflowServerConfig, message *bytes.Buffer, redisUrl string) {
	storer := RedisStorer{
		Address:  redisUrl,
		Password: "",
		Database: 0,
	}

	storer.Create()
	defer storer.Close()
	storer.Send(config.Owner, config.Repo, message)
}
