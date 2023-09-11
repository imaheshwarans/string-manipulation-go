package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"stringinator-go/config"
	"stringinator-go/router"
	"time"
)

type App struct {
	config    config.Configuration
	LogWriter *log.Logger
}

// func (app *App) logWriter() io.Writer {
// 	if app.LogWriter != nil {
// 		return app.LogWriter
// 	}
// 	return os.Stderr
// }

// func (app *App) configureLogs(stdOut, logFile bool) error {
// 	var ioWriterDefault io.Writer
// 	ioWriterDefault = app.logWriter()
// 	if stdOut {
// 		if logFile {
// 			ioWriterDefault = io.MultiWriter(os.Stdout, app.logWriter())
// 		} else {
// 			ioWriterDefault = os.Stdout
// 		}
// 	}

// 	lv, err := logrus.ParseLevel(app.config.LogLevel)
// 	if err != nil {
// 		return errors.Wrap(err, "Failed to initiate loggers. Invalid log level: "+app.config.LogLevel)
// 	}
// 	// f := commLog.LogFormatter{MaxLength: logConfig.MaxLength}
// 	// commLogInt.SetLogger(commLog.DefaultLoggerName, lv, &f, ioWriterDefault, false)
// 	// commLogInt.SetLogger(commLog.SecurityLoggerName, lv, &f, ioWriterSecurity, false)

// 	// secLog.Info(commLogMsg.LogInit)
// 	// defaultLog.Info(commLogMsg.LogInit)
// 	return nil
// }

func (a *App) startServer() {

	defaultLog.Println("server.go:Entering startServer()")
	defer defaultLog.Println("server.go:Leaving startServer()")

	router := router.InitRoutes()
	// stop := make(chan os.Signal)
	// go func() {
	err := http.ListenAndServe(":8888", router)
	if err != nil {
		fmt.Println("Failed to start server")
		// stop <- syscall.SIGTERM
	}
	// }()

	// <-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
