package utils

import (
	"fmt"
	"os"
	"stringinator-go/constants"

	log "github.com/sirupsen/logrus"
)

func openLogFiles() (logFile *os.File, err error) {

	logFile, err = os.OpenFile(constants.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(constants.LogFile, 0640)
	if err != nil {
		return nil, err
	}
	return
}

func ConfigureLogs(level string) *log.Logger {
	logFile, err := openLogFiles()
	if err != nil {
		fmt.Println("Failed to open log file")
	}

	log.SetOutput(logFile)
	l := log.New()
	lv, err := log.ParseLevel(level)
	if err != nil {
		fmt.Println("Failed to parse log level")
	}
	l.SetLevel(lv)
	return l
}
