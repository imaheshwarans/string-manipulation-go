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
	config    *config.Configuration
	LogWriter *log.Logger
}

func (a *App) startServer() {

	defaultLog.Trace("server.go:Entering startServer()")
	defer defaultLog.Trace("server.go:Leaving startServer()")

	router := router.InitRoutes()

	err := http.ListenAndServe(":1323", router)
	if err != nil {
		fmt.Println("Failed to start server")
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
