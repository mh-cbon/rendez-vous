package admin

import (
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"
	"github.com/mh-cbon/rendez-vous/browser"
)

func New(srvListen, dbFile, static string, proxyApp *browser.Proxy) *Website {
	return &Website{
		srvListen: srvListen,
		dbFile:    dbFile,
		static:    static,
		proxyApp:  proxyApp,
	}
}

type Website struct {
	srv       *http.Server
	srvListen string
	db        *storm.DB
	dbFile    string
	static    string
	proxyApp  *browser.Proxy
}

func (w *Website) ListenAndServe() error {
	db, err := storm.Open(w.dbFile)
	if err != nil {
		return err
	}
	w.db = db

	r := mux.NewRouter()
	if w.proxyApp != nil {
		Node(r, w.proxyApp)
	}
	Apps(r, w.db)
	Static(r, w.static)

	w.srv = &http.Server{
		Handler:      r,
		Addr:         w.srvListen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return w.srv.ListenAndServe()
}

func (w *Website) Close() error {

	err := w.db.Close()
	if err == nil {
		w.db = nil
	}

	err = w.srv.Close()
	if err == nil {
		w.srv = nil
	}

	return err
}
