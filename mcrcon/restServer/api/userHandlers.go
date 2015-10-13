// Handle the /api/users routes

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joshproehl/minecontrol/mcrcon"
	"net/http"
)

// Handle a request to the /users resource
func UsersRootHandler(w http.ResponseWriter, r *http.Request) {

	user1 := mcrcon.User{Username: "joshproehl"}
	user2 := mcrcon.User{Username: "flyingdwarves"}

	userList := mcrcon.Users{user1, user2}

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		panic(err)
	}
}

// Handle a request tho the root reesource
func UsernameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if err := json.NewEncoder(w).Encode(username); err != nil {
		panic(err)
	}
}
