package html

import (
	"net/http"
)

//
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HTML TEMPLATE"))
}
