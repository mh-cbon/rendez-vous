package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/gorilla/mux"
)

type App struct {
	ID                 int
	Group              string `storm:"index"`
	Type               string `storm:"unique"`
	URL                string
	StartCommand       string
	KillCommand        string
	ExtraKill          bool
	Announce           bool
	RequireCredentials bool
	Status             string
	LastLogMessage     string
}

type addRequest struct {
	Group string
	Type  string
	URL   string
}
type updateRequest struct {
	ID                 int
	Group              string `storm:"index"`
	Type               string `storm:"unique"`
	URL                string
	StartCommand       string
	KillCommand        string
	ExtraKill          bool
	Announce           bool
	RequireCredentials bool
}

var (
	GoApp = "GO"
)

var (
	StatusAnnouncing = "Announcing"
	StatusDisabled   = "Disabled"
	StatusEnabled    = "Enabled"
)

var okTypes = []string{GoApp}

func Apps(r *mux.Router, db *storm.DB) {

	r.HandleFunc("/list/{start:[0-9]+}/{limit:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var start int
		var limit int
		var res []App
		err := dontFail(w, jsonEncode(&res),
			urlIntVar(r, "start", &start),
			urlIntVar(r, "limit", &limit),
			dbReadApps(db, &res, &start, &limit),
		)
		log.Println(err)

	}).Methods("GET")

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var req addRequest
		err := dontFail(w, jsonEncode(&req),
			jsonDecode(r, &req),
			notEmpty(&req.URL, "URL"),
			notEmpty(&req.Type, "Type"),
			enum(okTypes, &req.Type, "Type"),
			dbSave(db, &req),
		)
		log.Println(err)
	}).Methods("POST")

	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		var req updateRequest
		var err error
		var tx storm.Node
		defer txFinish(tx, err)
		err = dontFail(w, jsonEncode(&req),
			txAcquire(db, &tx),
			jsonDecode(r, &req),
			notEmpty(&req.URL, "URL"),
			uniqueAppURL(tx, &req.URL, &req.ID),
			notEmpty(&req.Type, "Type"),
			enum(okTypes, &req.Type, "Type"),
			dbSave(tx, &req),
		)
		log.Println(err)
	}).Methods("POST")

	r.HandleFunc("/status/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var res App
		var id int
		err := dontFail(w, jsonEncode(&res),
			urlIntVar(r, "id", &id),
			dbReadAppByID(db, &res, &id),
		)
		log.Println(err)
	}).Methods("POST")

	r.HandleFunc("/start/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var res App
		var id int
		var err error
		var tx storm.Node
		defer txFinish(tx, err)
		err = dontFail(w, jsonEncode(&res),
			txAcquire(db, &tx),
			urlIntVar(r, "id", &id),
			dbReadAppByID(db, &res, &id),
			appStart(db, &res),
			dbSave(db, &res),
		)
		log.Println(err)
	}).Methods("POST")

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		var res App
		var id int
		var err error
		var tx storm.Node
		defer txFinish(tx, err)
		err = dontFail(w, jsonEncode(&res),
			txAcquire(db, &tx),
			urlIntVar(r, "id", &id),
			dbReadAppByID(tx, &res, &id),
			appStop(tx, &res),
			dbSave(tx, &res),
		)
		log.Println(err)
	}).Methods("POST")

}

func txFinish(tx storm.Tx, err error) {
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}

func txAcquire(db *storm.DB, tx *storm.Node) func() error {
	return func() error {
		x, err := db.Begin(true)
		if err != nil {
			return err
		}
		tx = &x
		return nil
	}
}

func appStart(db storm.Node, app *App) func() error {
	return func() error {
		if app == nil {
			return fmt.Errorf("app is nil")
		}
		if app.Status == StatusAnnouncing {
			return fmt.Errorf("app is started")
		}
		if app.Status == StatusDisabled {
			return fmt.Errorf("app is disabled")
		}
		return fmt.Errorf("I dont know what to do yet.")
	}
}

func appStop(db storm.Node, app *App) func() error {
	return func() error {
		if app == nil {
			return fmt.Errorf("app is nil")
		}
		if app.Status != StatusAnnouncing {
			return fmt.Errorf("app is not started")
		}
		return fmt.Errorf("I dont know what to do yet.")
	}
}

func dbReadAppByID(db storm.Node, data *App, id *int) func() error {
	return func() error {
		if id == nil {
			return fmt.Errorf("id must not be nil")
		}
		n := *id
		query := db.Select(q.Eq("ID", n)).Limit(1)
		return query.Find(data)
	}
}

func dbReadApps(db storm.Node, data *[]App, start, limit *int) func() error {
	return func() error {
		if start == nil {
			return fmt.Errorf("start must not be nil")
		}
		if limit == nil {
			return fmt.Errorf("limit must not be nil")
		}
		s := *start
		l := *limit
		query := db.Select().Skip(s).Limit(l)
		return query.Find(data)
	}
}

func urlIntVar(r *http.Request, name string, n *int) func() error {
	vars := mux.Vars(r)
	return func() error {
		if x, ok := vars[name]; ok {
			u, err := strconv.Atoi(x)
			if err == nil {
				n = &u
			}
			return err
		}
		return fmt.Errorf("url parmeter %q not found", name)
	}
}

func jsonEncode(data interface{}) func(http.ResponseWriter) error {
	return func(w http.ResponseWriter) error {
		return json.NewEncoder(w).Encode(data)
	}
}

func jsonDecode(r *http.Request, data interface{}) func() error {
	return func() error {
		return json.NewDecoder(r.Body).Decode(data)
	}
}

func dbSave(db storm.Node, data interface{}) func() error {
	return func() error {
		return db.Save(data)
	}
}

func uniqueAppURL(db storm.Node, url *string, id *int) func() error {
	return func() error {
		if url == nil {
			return fmt.Errorf("url is nil")
		}
		if id == nil {
			return fmt.Errorf("url is nil")
		}
		u := *url
		i := *id
		var similars []*App
		query := db.Select(q.Eq("URL", u), q.Not(q.Eq("ID", i))).Limit(1)
		c, err := query.Count(similars)
		if err != nil {
			return err
		}
		if c > 0 {
			return fmt.Errorf("similar app url name found: %v", similars[0].URL)
		}
		return nil
	}
}

func notEmpty(s *string, args ...interface{}) func() error {
	return func() error {
		if s == nil {
			return fmt.Errorf("must not be nil:%v", args...)
		}
		v := *s
		if v == "" {
			return fmt.Errorf("must not be nil:%v", args...)
		}
		return nil
	}
}

func enum(enum []string, s *string, args ...interface{}) func() error {
	return func() error {
		if s == nil {
			return fmt.Errorf("must not be nil:%v", args...)
		}
		v := *s
		for _, e := range enum {
			if v == e {
				return nil
			}
		}
		args = append([]interface{}{v, enum}, args...)
		return fmt.Errorf("%v must be in %v:%v", args...)
	}
}

func dontFail(w http.ResponseWriter, success func(http.ResponseWriter) error, h ...func() error) error {
	for _, hh := range h {
		if err := hh(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return success(w)
}
