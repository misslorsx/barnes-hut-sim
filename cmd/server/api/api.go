package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tardisman5197/barnes-hut-sim/pkg/simulation"
)

// API provides a server that handels request related
// to a barnes hut simulation.
type API struct {
	server *http.Server

	simulations map[string]simulation.Simulation
}

// NewAPI returns an instance of an API struct.
func NewAPI() API {
	var a API
	a.setup()
	a.simulations = make(map[string]simulation.Simulation)
	return a
}

// setup creates the http server.
func (a *API) setup() {
	r := mux.NewRouter()
	r.HandleFunc("/simulation/new", a.newSimulation).Methods("POST")
	r.HandleFunc("/simulation/start/{simID}/{steps}", a.start).Methods("GET")
	r.HandleFunc("/simulation/status/{simID}", a.status).Methods("GET")
	r.HandleFunc("/simulation/results/{simID}", a.results).Methods("GET")
	r.HandleFunc("/simulation/remove/{simID}", a.remove).Methods("GET")

	a.server = &http.Server{
		Handler:      r,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// Listen starts the api server listening for
// requests.
func (a *API) Listen() chan bool {
	done := make(chan bool, 1)

	go func() {
		err := a.server.ListenAndServe()
		if err != nil {
			done <- true
			return
		}
	}()

	return done
}

// Shutdown gracefully stops the api server.
func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
