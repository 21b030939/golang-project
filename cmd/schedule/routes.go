package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	schedule1 := r.PathPrefix("/api/v1").Subrouter()

	// schedule Singleton
	// localhost:8081/api/v1/schedules
	schedule1.HandleFunc("/schedules", app.getScheduleList).Methods("GET")
	// Create a new schedule
	schedule1.HandleFunc("/schedules", app.createScheduleHandler).Methods("POST")
	// Get a specific schedule
	schedule1.HandleFunc("/schedules/{id:[0-9]+}", app.getScheduleHandler).Methods("GET")
	// Update a specific schedule
	schedule1.HandleFunc("/schedules/{id:[0-9]+}", app.updateScheduleHandler).Methods("PUT")
	// Delete a specific schedule
	schedule1.HandleFunc("/schedules/{id:[0-9]+}", app.requirePermissions("schedules:write", app.deleteScheduleHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
