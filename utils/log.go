package utils

import (
	"fmt"
	"log"
	"os"
	"stringinator-go/constants"
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

func ConfigureLogs() *log.Logger {
	logFile, err := openLogFiles()
	if err != nil {
		fmt.Println("Failed to open log file")
	}

	log.SetOutput(logFile)
	return log.Default()
}
