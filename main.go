package main

import (
	"fmt"
	"log"

	"stringinator-go/config"
	"stringinator-go/constants"
	"stringinator-go/utils"

	"github.com/sirupsen/logrus"
)

var defaultLog *logrus.Logger

func main() {
	fmt.Println("Starting Server")
	defer fmt.Println("Server Ended")

	config, err := config.LoadConfiguration()
	if err != nil {
		fmt.Println("Failed to load configuration ", err)
		return
	}

	err = config.Save(constants.ConfigFile)
	if err != nil {
		fmt.Println("Failed to save configuration ", err)
		return
	}

	defaultLog = utils.ConfigureLogs(config.LogLevel)

	defaultLog.Info("Log initiated")
	defaultLog.Trace("Log initiated")

	app := &App{
		config:    config,
		LogWriter: log.Default(),
	}
	app.startServer()
}
