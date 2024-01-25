package lib2

import (
	log "github.com/sirupsen/logrus"
)

func LogOnError(err error, msg string, level string) {
	if err != nil {
		if level == "" {
			level = "debug"
		}
		if level == "debug" {
			log.Debugf("%s: %s", msg, err)
		}
		if level == "warn" {
			log.Warnf("%s: %s", msg, err)
		}
		if level == "panic" {
			log.Panicf("%s: %s", msg, err)
		}
	}
}
