package restServer

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-zoo/bone"
	"github.com/joshproehl/minecontrol/mcrcon"
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

var rcon_client *mcrcon.MCRCONClient

// By default go generate is going to build the production version. Run the command with -debug flag for
// easier local development of static assets.
//go:generate go-bindata-assetfs -pkg restServer -prefix "gui/assets/" gui/assets/...

// NewServer creates a server that will listen for requests over HTTP and interact with the RCON server specified
// non-/api prefixed routes are served from static files compiled into bindata_assetfs.go
func NewRestServer(c *ServerConfig) {
	var err error
	rcon_client, err = mcrcon.NewClient(c.RCON_address, c.RCON_port, c.RCON_password)

	if err != nil {
		panic(fmt.Errorf("Could not connect to RCON server at %s:%d. (Error was: %s)", c.RCON_address, c.RCON_port, err))
	}

	router := bone.New()

	// Redirect static resources, and then handle the static resources (/gui/) routes with the static asset file
	router.Handle("/", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/gui/", 302)
	}))
	router.Get("/gui/", http.StripPrefix("/gui/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""})))

	// Define the API (JSON) routes
	router.GetFunc("/api", apiRootHandler)
	router.GetFunc("/api/users", usersRootHandler)
	router.GetFunc("/api/users/:username", usernameHandler)

	// TODO: Require a http basic auth username and password if passed in.

	// Start the server
	fmt.Println("Starting server on port", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), router)
}
