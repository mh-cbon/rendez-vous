package admin

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gocraft/dbr"
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
	conn      *dbr.Connection
	sess      *dbr.Session
	dbFile    string
	static    string
	proxyApp  *browser.Proxy
	tmpDir    string
}

func (w *Website) ListenAndServe() error {

	tmpDir, err := ioutil.TempDir("", "rendez-vous")
	if err != nil {
		return err
	}
	w.tmpDir = tmpDir

	conn, err := dbr.Open("sqlite3", w.dbFile, nil)
	if err != nil {
		return err
	}
	w.conn = conn
	w.sess = conn.NewSession(nil)

	r := mux.NewRouter()
	if w.proxyApp != nil {
		Node(r, w.proxyApp)
	}
	Apps(r, w.sess, w.tmpDir)
	UI(r, w.static)

	w.srv = &http.Server{
		Handler:      r,
		Addr:         w.srvListen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return w.srv.ListenAndServe()
}

func (w *Website) Close() error {

	err := os.RemoveAll(w.tmpDir)

	err = w.conn.Close()
	if err == nil {
		w.conn = nil
	}

	err = w.srv.Close()
	if err == nil {
		w.srv = nil
	}

	return err
}
