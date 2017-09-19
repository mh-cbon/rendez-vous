package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocraft/dbr"
	"github.com/gorilla/mux"
	"github.com/mh-cbon/rendez-vous/admin/shellexec"
	"github.com/mh-cbon/rendez-vous/cli"
)

type App struct {
	ID                 int64  `db:"id"`
	Type               string `db:"type"`
	URL                string `db:"url"`
	Name               string `db:"name"`
	BinaryName         string `db:"binary_name"`
	StartPattern       string `db:"start_pattern"`
	StartParams        []cli.Param
	KillPattern        string `db:"kill_pattern"`
	KillParams         []cli.Param
	StartCommand       string    `db:"start_command"`
	KillCommand        string    `db:"kill_command"`
	IsSystem           bool      `db:"is_system"`
	ExtraKill          bool      `db:"extra_kill"`
	AnnounceName       string    `db:"announce_name"`
	Announce           bool      `db:"announce"`
	RequireCredentials bool      `db:"require_credentials"`
	Credentials        string    `db:"credentials"`
	IsInstalled        bool      `db:"is_installed"`
	PassTest           bool      `db:"pass_test"`
	Status             string    `db:"status"`
	LastLogMessage     string    `db:"lastlog"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func (a *App) AutomaticName() string {
	if a.Name == "" {
		return filepath.Base(a.URL)
	}
	return a.Name
}

func execTemplate(tpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "", data)
	return b.String(), err
}

func (a *App) GenerateCommands(startData map[string]interface{}, killData map[string]interface{}) error {
	if a.StartCommand == "" {
		t, err := execTemplate(a.StartPattern, startData)
		if err != nil {
			return err
		}
		a.StartCommand = t
		if a.StartCommand == "" {
			return fmt.Errorf("start command is empty")
		}
	}
	if a.KillCommand == "" && a.KillPattern != "" {
		t, err := execTemplate(a.KillPattern, killData)
		if err != nil {
			return err
		}
		a.KillCommand = t
	}
	return nil
}

func installApp(app *App, out io.Writer) error {
	var cmdStr string
	if app.Type == GoApp {
		cmdStr = "go get -u -x  " + app.URL
	} else if app.Type == NpmApp {
		cmdStr = "npm i " + app.URL
	}
	out.Write([]byte(cmdStr))
	cmd, err := shellexec.Command("", cmdStr)
	if err != nil {
		return err
	}
	res, err := cmd.CombinedOutput()
	if res != nil {
		out.Write(res)
	}
	return err
}

func getConfig(app *App, out io.Writer) (*cli.App, error) {
	name := app.AutomaticName()
	cmdStr := name + " rendez-vous"
	out.Write([]byte(cmdStr))
	cmd, err := shellexec.Command("", cmdStr)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	res, err := cmd.CombinedOutput()
	if res != nil {
		out.Write(res)
	}
	if err != nil {
		return nil, err
	}
	buf.Write(res)
	cli := new(cli.App)
	if err := json.NewDecoder(buf).Decode(cli); err != nil {
		return nil, err
	}
	return cli, nil
}

var (
	GoApp  = "GO"
	NpmApp = "NPM"
)

var (
	StatusDisabled = "Disabled"
	StatusRunning  = "Running"
	StatusStopped  = "Stopped"
)

var okTypes = []string{GoApp, NpmApp}

func Apps(r *mux.Router, db *dbr.Session, tmpDir string) {

	appsModel := AppsModel{}

	apps := appsModel.FromDb(db)
	fail(reduce(
		func() error { return appsModel.Drop(db) },
		func() error { return appsModel.Setup(db) },
		func() error { return apps.Insert(&App{Name: "Admin", IsSystem: true, Status: StatusRunning}) },
		func() error { return apps.Insert(&App{Name: "Web", IsSystem: true, Status: StatusRunning}) },
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
		<-time.After(time.Second)
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
				func() error { return txApps.AddUserApp(req) },
				func() error { return json.NewEncoder(w).Encode(req) },
			)
		})
		<-time.After(time.Second)
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/install/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		app := new(App)
		var id int64
		vars := mux.Vars(r)
		err := tx(db, func(tx dbr.SessionRunner) error {
			txApps := appsModel.FromDb(tx)
			return reduce(
				func() error { return intVars(vars, "id", &id) },
				func() error { return apps.ByID(id, app) },
				func() error { return notEmptyString(app.URL, "URL") },
				func() error { return notEmptyString(app.Type, "Type") },
				func() error { return inEnum(okTypes, app.Type, "Type") },
				func() error {
					f := fmt.Sprintf("%v/install_%v.log", tmpDir, app.ID)
					return writeFile(f, 0644, func(logFile *os.File) error {
						out := io.MultiWriter(logFile, os.Stderr)
						return installApp(app, out)
					})
				},
				func() error {
					return txApps.SetAppInstalled(app)
				},
				func() error { return json.NewEncoder(w).Encode(app) },
			)
		})
		httpFail(w, err)
	}).Methods("POST")

	r.HandleFunc("/logs/{type:[a-z]+}/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		app := new(App)
		var typeStr string
		var id int64
		vars := mux.Vars(r)
		err := reduce(
			func() error { return strVars(vars, "type", &typeStr) },
			func() error { return intVars(vars, "id", &id) },
			func() error { return apps.ByID(id, app) },
			func() error {
				f := fmt.Sprintf("%v/%v_%v.log", tmpDir, typeStr, app.ID)
				return readFile(f, copyTo(w))
			},
		)
		httpFail(w, err)
	})

	r.HandleFunc("/config/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		app := new(App)
		var id int64
		vars := mux.Vars(r)
		err := tx(db, func(tx dbr.SessionRunner) error {
			txApps := appsModel.FromDb(tx)
			return reduce(
				func() error { return intVars(vars, "id", &id) },
				func() error { return apps.ByID(id, app) },
				func() error { return notEmptyString(app.URL, "URL") },
				func() error { return notEmptyString(app.Type, "Type") },
				func() error { return inEnum(okTypes, app.Type, "Type") },
				func() error {
					f := fmt.Sprintf("%v/config_%v.log", tmpDir, app.ID)
					return writeFile(f, 0644, func(logFile *os.File) error {
						out := io.MultiWriter(logFile, os.Stderr)
						cfg, err := getConfig(app, out)
						if err != nil {
							return err
						}
						app.StartPattern = cfg.Start
						app.KillPattern = cfg.Kill
						app.KillParams = cfg.KillParams
						app.StartParams = cfg.StartParams
						return app.GenerateCommands(cfg.StartValues(), cfg.KillValues())
					})
				},
				func() error {
					return txApps.AppTestPassed(app)
				},
				func() error { return json.NewEncoder(w).Encode(app) },
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
				func() error { return txApps.UpdateUserApp(req) },
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
				func() error { return txApps.DeleteUserApp(id) },
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
		<-time.After(time.Second)
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

func copyTo(dst io.Writer) func(*os.File) error {
	return func(src *os.File) error {
		_, err := io.Copy(dst, src)
		return err
	}
}

func openFile(filepath string, flags int, perm os.FileMode, h func(*os.File) error) error {
	f, err := os.OpenFile(filepath, flags, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	return h(f)
}

func appendFile(filepath string, perm os.FileMode, h func(*os.File) error) error {
	return openFile(filepath, os.O_CREATE|os.O_WRONLY, perm, h)
}

func writeFile(filepath string, perm os.FileMode, h func(*os.File) error) error {
	return openFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm, h)
}

func readFile(filepath string, h func(*os.File) error) error {
	return openFile(filepath, os.O_RDONLY, 0, h)
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

func strVars(vars map[string]string, name string, out *string) error {
	if x, ok := vars[name]; ok {
		*out = x
		return nil
	}
	return fmt.Errorf("not found %q", name)
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
