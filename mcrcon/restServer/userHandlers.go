// Handle the /api/users routes

package restServer

import (
	"encoding/json"
	"github.com/go-zoo/bone"
	"net/http"
)

// Handle a request to the /users resource
func usersRootHandler(w http.ResponseWriter, r *http.Request) {
	userList, cmdErr := rcon_client.SendCommand("/list")

	if cmdErr != nil {
		panic(cmdErr)
	}

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		panic(err)
	}
}

// Handle a request tho the root reesource
func usernameHandler(w http.ResponseWriter, r *http.Request) {
	username := bone.GetValue(r, "username")

	if err := json.NewEncoder(w).Encode(username); err != nil {
		panic(err)
	}
}
