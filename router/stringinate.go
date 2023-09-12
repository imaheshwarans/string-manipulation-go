package router

import (
	"stringinator-go/controllers"

	"github.com/gorilla/mux"
)

func setStringRoutes(router *mux.Router) *mux.Router {
	stringController := controllers.StringInate{}

	router.Handle("/stringinate", ErrorHandler(JsonResponseHandler(stringController.Create))).Methods("POST")
	router.Handle("/stringinate", ErrorHandler(JsonResponseHandler(stringController.Get))).Methods("GET")
	router.Handle("/stringinates", ErrorHandler(JsonResponseHandler(stringController.GetAll))).Methods("GET")
	return router
}
