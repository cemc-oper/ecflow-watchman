package ecflow_watchman

import (
	"github.com/nwpc-oper/ecflow-client-go"
	"github.com/pquerna/ffjson/ffjson"
	log "github.com/sirupsen/logrus"
	"time"
)

type EcflowServerStatus struct {
	StatusRecords []ecflow_client.StatusRecord `json:"status_records"`
	CollectedTime time.Time                    `json:"collected_time"`
}

func GetEcflowStatus(config EcflowServerConfig) *EcflowServerStatus {
	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Infof("get nodes...")

	client := ecflow_client.CreateEcflowClient(config.Host, config.Port)
	client.SetConnectTimeout(config.ConnectTimeout)
	defer client.Close()

	ret := client.Sync()
	if ret != 0 {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Errorf("sync has error: %v", ret)
		return nil
	}

	recordsJson := client.StatusRecordsJson()

	var records []ecflow_client.StatusRecord
	err := ffjson.Unmarshal([]byte(recordsJson), &records)

	if err != nil {
		log.WithFields(log.Fields{
			"owner": config.Owner,
			"repo":  config.Repo,
		}).Errorf("Unmarshal recordsJson has error: %v", err)
		return nil
	}

	ecflowServerStatus := &EcflowServerStatus{
		StatusRecords: records,
		CollectedTime: client.CollectedTime,
	}

	log.WithFields(log.Fields{
		"owner": config.Owner,
		"repo":  config.Repo,
	}).Infof("get nodes...done")

	return ecflowServerStatus
}
