package router

import (
	"stringinator-go/controllers"

	"github.com/gorilla/mux"
)

func setStatsRoutes(router *mux.Router) *mux.Router {
	stateController := controllers.Stats{}

	router.Handle("/stats", ErrorHandler(JsonResponseHandler(stateController.GetStats))).Methods("GET")
	return router
}
