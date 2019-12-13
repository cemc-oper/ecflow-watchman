package main

import (
	"github.com/nwpc-oper/ecflow-watchman/cmd/cli"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
		FullTimestamp:   true,
	})
}

func main() {
	cli.Execute()
}
