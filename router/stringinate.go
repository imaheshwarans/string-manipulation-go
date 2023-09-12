package router

import (
	"stringinator-go/controllers"
	"stringinator-go/utils"

	"github.com/gorilla/mux"
)

func setStringRoutes(router *mux.Router) *mux.Router {
	stringController := controllers.StringInate{}

	router.Handle("/stringinate", utils.ErrorHandler(utils.JsonResponseHandler(stringController.Create))).Methods("POST")
	router.Handle("/stringinate", utils.ErrorHandler(utils.JsonResponseHandler(stringController.Get))).Methods("GET")
	router.Handle("/stringinates", utils.ErrorHandler(utils.JsonResponseHandler(stringController.GetAll))).Methods("GET")
	return router
}
