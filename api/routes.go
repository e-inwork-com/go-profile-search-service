package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(
		http.MethodGet, "/service/search/profiles",
		app.getProfileHandler)

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
