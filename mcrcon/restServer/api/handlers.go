// Handlers and tools for the /api/ routes of MCRCON's server
package api

import (
	"net/http"
)

// Handle a request tho the root reesource
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("JSON DATA"))
}
