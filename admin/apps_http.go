package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocraft/dbr"
	"github.com/gorilla/mux"
)

type App struct {
	ID                 int64     `db:"id"`
	Type               string    `db:"type"`
	URL                string    `db:"url"`
	Name               string    `db:"name"`
	StartCommand       string    `db:"start_command"`
	KillCommand        string    `db:"kill_command"`
	IsSystem           bool      `db:"is_system"`
	ExtraKill          bool      `db:"extra_kill"`
	Announce           bool      `db:"announce"`
	RequireCredentials bool      `db:"require_credentials"`
	IsInstalled        bool      `db:"is_installed"`
	Status             string    `db:"status"`
	LastLogMessage     string    `db:"lastlog"`
	UpdatedAt          time.Time `db:"updated_at"`
}

var (
	GoApp  = "GO"
	NpmApp = "GO"
)

var (
	StatusAnnouncing = "Announcing"
	StatusDisabled   = "Disabled"
	StatusEnabled    = "Enabled"
)

var okTypes = []string{GoApp, NpmApp}

func Apps(r *mux.Router, db *dbr.Session) {

	appsModel := AppsModel{}

	apps := appsModel.FromDb(db)
	fail(reduce(
		func() error { return appsModel.Drop(db) },
		func() error { return appsModel.Setup(db) },
		func() error { return apps.Insert(&App{Name: "Admin", IsSystem: true}) },
		func() error { return apps.Insert(&App{Name: "Web", IsSystem: true}) },
	))

	r.HandleFunc("/list/{start:[0-9]+}/{limit:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var start uint64
		var limit uint64
		res := []*App{}
		err := reduce(
			func() error { return uintVars(vars, "start", &start) },
			func() error { return uintVars(vars, "limit", &limit) },
			func() error { return apps.Many(start, limit, &res) },
			func() error { return json.NewEncoder(w).Encode(res) },
		)
		httpFail(w, err)
	}).Methods("GET")

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		req := new(App)
		err := tx(db, func(tx dbr.SessionRunner) error {
			txApps := appsModel.FromDb(tx)
			return reduce(
				func() error { return json.NewDecoder(r.Body).Decode(&req) },
				//todo: stronger validation
				func() error { return notEmptyString(req.URL, "URL") },
				func() error { return notEmptyString(req.Type, "Type") },
				func() error { return inEnum(okTypes, req.Type, "Type") },
				func() error { return txApps.IsURLUnique(req.URL, req.Type, nil) },
				func() error { return txApps.Insert(req) },
				//todo: download + install app
				func() error { return json.NewEncoder(w).Encode(req) },
			)
		})
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		req := new(App)
		err := tx(db, func(tx dbr.SessionRunner) error {
			txApps := appsModel.FromDb(tx)
			return reduce(
				func() error { return json.NewDecoder(r.Body).Decode(&req) },
				//todo: stronger validation
				func() error { return notEmptyString(req.URL, "URL") },
				func() error { return notEmptyString(req.Type, "Type") },
				func() error { return inEnum(okTypes, req.Type, "Type") },
				func() error { return txApps.IsURLUnique(req.URL, req.Type, &req.ID) },
				func() error { return txApps.Update(req) },
				func() error { return json.NewEncoder(w).Encode(req) },
			)
		})
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/delete/{id:[0-9]}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data := new(App)
		var id int64
		err := tx(db, func(tx dbr.SessionRunner) error {
			txApps := appsModel.FromDb(tx)
			return reduce(
				func() error { return intVars(vars, "id", &id) },
				func() error { return txApps.ByID(id, data) },
				func() error { return notSystemApp(data) },
				func() error { return txApps.DeleteByID(id) },
				func() error { return json.NewEncoder(w).Encode(data) },
			)
		})
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/status/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		res := new(App)
		var id int64
		vars := mux.Vars(r)
		err := reduce(
			func() error { return intVars(vars, "id", &id) },
			func() error { return apps.ByID(id, res) },
			func() error { return json.NewEncoder(w).Encode(res) },
		)
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/start/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var res *App
		var id int64
		vars := mux.Vars(r)
		err := reduce(
			func() error { return intVars(vars, "id", &id) },
			// todo: start app
			func() error { return json.NewEncoder(w).Encode(res) },
		)
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		var res *App
		var id int64
		vars := mux.Vars(r)
		err := reduce(
			func() error { return intVars(vars, "id", &id) },
			// todo: stop app
			func() error { return json.NewEncoder(w).Encode(res) },
		)
		httpFail(w, err)
	}).Methods("POST")

}

func tx(db *dbr.Session, h func(tx dbr.SessionRunner) error) error {
	x, err := db.Begin()
	if err != nil {
		return err
	}
	if err := h(x); err != nil {
		x.Rollback()
		return err
	}
	return x.Commit()
}

func intVars(vars map[string]string, name string, out *int64) error {
	if x, ok := vars[name]; ok {
		y, err := strconv.ParseInt(x, 10, 0)
		if err == nil {
			*out = y
		}
		return err
	}
	return fmt.Errorf("not found %q", name)
}

func uintVars(vars map[string]string, name string, out *uint64) error {
	if x, ok := vars[name]; ok {
		y, err := strconv.ParseUint(x, 10, 0)
		if err == nil {
			*out = y
		}
		return err
	}
	return fmt.Errorf("not found %q", name)
}

func inEnum(enum []string, s string, name string) error {
	for _, e := range enum {
		if s == e {
			return nil
		}
	}
	return fmt.Errorf("%v must be in %v:%v", name, enum, s)
}

func notEmptyString(s string, name string) error {
	if s == "" {
		return fmt.Errorf("must %q not be empty", name)
	}
	return nil
}

func notSystemApp(app *App) error {
	if app.IsSystem {
		return fmt.Errorf("cannot delete app %q, it is vital.", app.Name)
	}
	return nil
}

func reduce(h ...func() error) error {
	for _, hh := range h {
		if err := hh(); err != nil {
			return err
		}
	}
	return nil
}

func httpFail(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
