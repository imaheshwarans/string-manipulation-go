package router

import (
	"stringinator-go/controllers"
	"stringinator-go/utils"

	"github.com/gorilla/mux"
)

func setStatsRoutes(router *mux.Router) *mux.Router {
	stateController := controllers.Stats{}

	router.Handle("/stats", utils.ErrorHandler(utils.JsonResponseHandler(stateController.GetStats))).Methods("GET")
	return router
}
