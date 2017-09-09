package browser

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MakeWebsite for th me.com website
func MakeWebsite(proxyApp *Proxy, dir string) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/get_port", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Status int
			Port   int
		}{Port: proxyApp.Port().Num(), Status: proxyApp.Port().Status()}
		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(b)
	}).Methods("POST")

	r.HandleFunc("/test_port", func(w http.ResponseWriter, r *http.Request) {
		status := proxyApp.TestPort()
		res := struct {
			Status int
			Port   int
		}{Port: status.Num(), Status: status.Status()}
		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(b)
	}).Methods("POST")

	r.HandleFunc("/change_port/{port:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		port := vars["port"]

		if err := proxyApp.ChangeListenAddress(":" + port); err != nil {
			http.Error(w, "addr: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res := struct {
			Status int
			Port   int
		}{Port: proxyApp.Port().Num(), Status: proxyApp.Port().Status()}
		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(b)
	}).Methods("POST")

	r.HandleFunc("/list/{start:[0-9]+}/{limit:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		start := vars["start"]
		limit := vars["limit"]
		s, err := strconv.Atoi(start)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		l, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := proxyApp.List(s, l)
		if err != nil {
			http.Error(w, "addr: "+err.Error(), http.StatusInternalServerError)
			return
		}
		peers := []Peer{}
		for _, r := range res {
			peers = append(peers, Peer{Pbk: hex.EncodeToString(r.Pbk), Name: r.Value})
		}

		b, err := json.Marshal(peers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(b)
	}).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))
	return r
}

type Peer struct {
	Name string
	Pbk  string
}
