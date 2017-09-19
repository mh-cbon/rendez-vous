package admin

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocraft/dbr"
)

var appTable = "apps"
var appDropTable = "DROP TABLE apps;"

var appCreateTable = `
CREATE TABLE IF NOT EXISTS apps (
	id  INTEGER PRIMARY KEY AUTOINCREMENT,
	type text,
	binary_name text,
	name text,
	url text,
	start_pattern text,
	kill_pattern text,
	start_command text,
	kill_command text,
	is_system int,
	extra_kill int,
	announce_name text,
	announce int,
	credentials text,
	require_credentials int,
	is_installed int,
	pass_test int,
	status text,
	lastlog text,
	updated_at date
);
`

type AppsModel struct{}

func (a AppsModel) Setup(db *dbr.Session) error {
	_, err := db.Exec(appCreateTable)
	return err
}

func (a AppsModel) Drop(db *dbr.Session) error {
	_, err := db.Exec(appDropTable)
	return err
}

func (a AppsModel) FromDb(db dbr.SessionRunner) AppsQuerier {
	return AppsQuerier{db}
}

type AppsQuerier struct {
	db dbr.SessionRunner
}

func (a AppsQuerier) ByID(id int64, out *App) error {
	res, err := a.Select().Where("id = ?", id).Load(&out)
	if res == 0 && err == nil {
		err = fmt.Errorf("not found %v", id)
	}
	return err
}

func (a AppsQuerier) DeleteByID(id int64) error {
	res, err := a.db.DeleteFrom(appTable).Where("id = ?", id).Exec()
	return mustAffectRows(res, err)
}

func (a AppsQuerier) IsURLUnique(url, Type string, notID *int64) error {
	var c int
	q := a.Count().Where("url = ?", url).Where("type = ?", Type)
	if notID != nil {
		q = q.Where("id != ?", *notID)
	}
	_, err := q.Load(&c)
	if err != nil {
		return err
	}
	if c > 0 {
		return fmt.Errorf("url %q is not unique", url)
	}
	return nil
}

func (a AppsQuerier) Many(start, limit uint64, out *[]*App) error {
	_, err := a.Select().Offset(start).Limit(limit).LoadStructs(out)
	return err
}

func (a AppsQuerier) Select(what ...string) *dbr.SelectBuilder {
	if len(what) == 0 {
		what = append(what, "*")
	}
	return a.db.Select(what...).From(appTable)
}

func (a AppsQuerier) Count(what ...string) *dbr.SelectBuilder {
	if len(what) == 0 {
		what = append(what, "COUNT(*)")
	}
	return a.Select(what...)
}

func (a AppsQuerier) Update(data *App) error {
	t := time.Now()
	res, err := a.db.Update(appTable).
		Set("type", data.Type).
		Set("url", data.URL).
		Set("name", data.Name).
		Set("binary_name", data.BinaryName).
		Set("start_pattern", data.StartPattern).
		Set("kill_pattern", data.KillPattern).
		Set("start_command", data.StartCommand).
		Set("kill_command", data.KillCommand).
		Set("is_system", data.IsSystem).
		Set("extra_kill", data.ExtraKill).
		Set("announce_name", data.AnnounceName).
		Set("announce", data.Announce).
		Set("credentials", data.Credentials).
		Set("require_credentials", data.RequireCredentials).
		Set("is_installed", data.IsInstalled).
		Set("pass_test", data.PassTest).
		Set("status", data.Status).
		Set("lastlog", data.LastLogMessage).
		Set("updated_at", t).
		Where("id = ?", data.ID).
		Exec()
	if err == nil {
		data.UpdatedAt = t
	}
	return mustAffectRows(res, err)
}

func (a AppsQuerier) Insert(data *App) error {
	data.UpdatedAt = time.Now()
	res, err := a.db.InsertInto(appTable).Columns(
		"type",
		"url",
		"name",
		"binary_name",
		"start_pattern",
		"kill_pattern",
		"start_command",
		"kill_command",
		"is_system",
		"extra_kill",
		"announce_name",
		"announce",
		"credentials",
		"require_credentials",
		"is_installed",
		"pass_test",
		"status",
		"lastlog",
		"updated_at",
	).Record(data).Exec()
	if err == nil {
		id, err2 := res.LastInsertId()
		if err2 != nil {
			return err2
		}
		data.ID = id
	}
	return err
}

func (a AppsQuerier) AddUserApp(data *App) error {
	data.UpdatedAt = time.Now()
	data.IsSystem = false
	data.IsInstalled = false
	data.PassTest = false
	data.Name = ""
	data.KillCommand = ""
	data.KillPattern = ""
	data.StartCommand = ""
	data.StartPattern = ""
	data.Credentials = ""
	data.AnnounceName = ""
	data.Announce = false
	data.ExtraKill = false
	return a.Insert(data)
}

func (a AppsQuerier) SetAppInstalled(data *App) error {
	data.IsInstalled = true
	t := time.Now()
	res, err := a.db.Update(appTable).
		Set("is_installed", data.IsInstalled).
		Set("lastlog", data.LastLogMessage).
		Set("updated_at", t).
		Where("id = ?", data.ID).
		Exec()
	if err == nil {
		data.UpdatedAt = t
	}
	return mustAffectRows(res, err)
}

func (a AppsQuerier) AppTestPassed(data *App) error {
	data.PassTest = true
	t := time.Now()
	res, err := a.db.Update(appTable).
		Set("pass_test", data.PassTest).
		Set("lastlog", data.LastLogMessage).
		Set("updated_at", t).
		Where("id = ?", data.ID).
		Exec()
	if err == nil {
		data.UpdatedAt = t
	}
	return mustAffectRows(res, err)
}

func (a AppsQuerier) UpdateUserApp(data *App) error {
	t := time.Now()
	res, err := a.db.Update(appTable).
		Set("binary_name", data.BinaryName).
		Set("start_command", data.StartCommand).
		Set("kill_command", data.KillCommand).
		Set("extra_kill", data.ExtraKill).
		Set("announce_name", data.AnnounceName).
		Set("announce", data.Announce).
		Set("credentials", data.Credentials).
		Set("require_credentials", data.RequireCredentials).
		Set("updated_at", t).
		Where("id = ?", data.ID).
		Exec()
	if err == nil {
		data.UpdatedAt = t
	}
	return mustAffectRows(res, err)
}

func (a AppsQuerier) DeleteUserApp(id int64) error {
	res, err := a.db.DeleteFrom(appTable).Where("id = ?", id).Where("is_system = ?", false).Exec()
	return mustAffectRows(res, err)
}

func mustAffectRows(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
