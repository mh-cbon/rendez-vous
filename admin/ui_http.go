package admin

import (
	"net/http"

	"github.com/gorilla/mux"
)

func UI(r *mux.Router, dir string) {
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))
}
