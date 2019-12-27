package ecflow_watchman

import (
	"encoding/json"
	"github.com/nwpc-oper/ecflow-client-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type EcflowServerStatus struct {
	StatusRecords json.RawMessage `json:"status_records"`
	CollectedTime time.Time       `json:"collected_time"`
}

func GetEcflowStatus(config EcflowServerConfig, client *ecflow_client.EcflowClient) *EcflowServerStatus {
	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Infof("get nodes...")

	ret := client.Sync()
	if ret != 0 {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Errorf("sync has error: %v", ret)
		return nil
	}

	recordsJson := client.StatusRecordsJson()

	ecflowServerStatus := &EcflowServerStatus{
		StatusRecords: json.RawMessage(recordsJson),
		CollectedTime: client.CollectedTime,
	}

	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Infof("get nodes...done")

	return ecflowServerStatus
}
