package main

import (
	"log"
	"stringinator-go/config"
	"stringinator-go/utils"
)

var defaultLog *log.Logger

func init() {
	defaultLog = utils.ConfigureLogs()
}

func main() {
	defaultLog.Println("Starting Server")
	defer defaultLog.Println("Server Ended")

	app := &App{
		config: config.Configuration{
			LogLevel:    "info",
			LogFilePath: "string.log",
		},
		LogWriter: log.Default(),
	}

	app.startServer()
}
