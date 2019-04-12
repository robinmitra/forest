package main

import (
	"github.com/robinmitra/forest/cmd"
	log "github.com/sirupsen/logrus"
)

var VERSION = "0.2.0"

func main() {
	cmd.Execute(VERSION)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.WarnLevel) // Only log the warning severity or above.
}
