package router

import (
	"stringinator-go/controllers"

	"github.com/gorilla/mux"
)

func setStringRoutes(router *mux.Router) *mux.Router {
	// defaultLog.Trace("router/version:setVersionRoutes() Entering")
	// defer defaultLog.Trace("router/version:setVersionRoutes() Leaving")
	stringController := controllers.StringInate{}

	// router.Handle("/stringinate", ErrorHandler(ResponseHandler(stringController.GetVersion))).Methods("GET").Methods("POST")
	router.Handle("/stringinate", ErrorHandler(JsonResponseHandler(stringController.Create))).Methods("POST")
	router.Handle("/stringinate", ErrorHandler(JsonResponseHandler(stringController.Get))).Methods("GET")
	return router
}
