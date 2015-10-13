package restServer

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/joshproehl/minecontrol/mcrcon/restServer/api"
	"net/http"
)

type ServerConfig struct {
	RCON_address  string
	RCON_port     int
	RCON_password string
	Username      string
	Password      string
	Port          int
}

// By default go generate is going to build the production version. Run the command with -debug flag for
// easier local development of static assets.
//go:generate go-bindata-assetfs -pkg restServer -prefix "gui/assets/" gui/assets/...

// NewServer creates a server that will listen for requests over HTTP and interact with the RCON server specified
// non-/api prefixed routes are served from static files compiled into bindata_assetfs.go
func NewRestServer(c *ServerConfig) {
	router := mux.NewRouter()

	// Define the API (JSON) routes
	api_router := router.PathPrefix("/api").Subrouter()
	api_router.HandleFunc("/", api.RootHandler)
	api_router.HandleFunc("/users", api.UsersRootHandler).Name("users")
	api_router.HandleFunc("/users/{username}", api.UsernameHandler)

	// Redirect static resources, and then handle the static resources (/gui/) routes with the static asset file
	router.Methods("GET").Path("/").HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/gui/", 302)
	})
	router.PathPrefix("/gui/").Handler(http.StripPrefix("/gui/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""})))

	// TODO: Require a http basic auth username and password if passed in.

	// Start the server
	fmt.Println("Starting server on port", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), router)
}
