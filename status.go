package ecflow_watchman

import (
	"github.com/perillaroc/ecflow-client-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type EcflowServerStatus struct {
	StatusRecords []ecflow_client.StatusRecord `json:"status_records"`
	CollectedTime time.Time                    `json:"collected_time"`
}

func GetEcflowStatus(config EcflowServerConfig) *EcflowServerStatus {
	client := ecflow_client.CreateEcflowClient(config.Host, config.Port)
	client.SetConnectTimeout(config.ConnectTimeout)
	defer client.Close()

	ret := client.Sync()
	if ret != 0 {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Error("sync has error: ", ret)
		return nil
	}

	records := client.StatusRecords()

	ecflowServerStatus := &EcflowServerStatus{
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

	return ecflowServerStatus
}
