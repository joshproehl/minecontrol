// Handlers and tools for the /api/ routes of MCRCON's server
package restServer

import (
	"net/http"
)

// Handle a request tho the root reesource
func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("JSON DATA"))
}
