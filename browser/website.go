package browser

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MakeWebsite for th me.com website
func MakeWebsite(dir string) http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))
	return r
}
