package mcrcon

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joshproehl/minecontrol-go/mcrcon/api"
	"github.com/joshproehl/minecontrol-go/mcrcon/html"
	"net/http"
)

// NewServer creates a server that will listen for requests over HTTP and interact with the RCON server specified
func NewRestServer(p_rcon_address string, p_rcon_password string, p_rcon_port int, p_username string, p_password string, p_port int) {
	// Define the HTML routes
	router := mux.NewRouter()
	router.HandleFunc("/", html.RootHandler)

	// Define the API (JSON) routes
	api_router := router.PathPrefix("/api").Subrouter()
	api_router.HandleFunc("/", api.RootHandler)

	fmt.Println("Starting server on port", p_port)

	// TODO: Require a http basic auth username and password if passed in.
	http.ListenAndServe(fmt.Sprintf(":%d", p_port), router) // TODO: Use the port passed-in.

}
