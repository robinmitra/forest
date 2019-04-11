package main

import (
	"github.com/robinmitra/forest/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.WarnLevel) // Only log the warning severity or above.
}
