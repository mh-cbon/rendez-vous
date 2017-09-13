package admin

import (
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
	name text,
	url text,
	start_command text,
	kill_command text,
	is_system int,
	extra_kill int,
	announce int,
	require_credentials int,
	is_installed int,
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
	_, err := a.db.DeleteFrom(appTable).Where("id = ?", id).Exec()
	return err
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
	_, err := a.db.Update(appTable).
		Set("type", data.Type).
		Set("url", data.URL).
		Set("name", data.Name).
		Set("start_command", data.StartCommand).
		Set("kill_command", data.KillCommand).
		Set("is_system", data.IsSystem).
		Set("extra_kill", data.ExtraKill).
		Set("announce", data.Announce).
		Set("require_credentials", data.RequireCredentials).
		Set("is_installed", data.IsInstalled).
		Set("status", data.Status).
		Set("lastlog", data.LastLogMessage).
		Set("updated_at", t).
		Where("id = ?", data.ID).
		Exec()
	if err == nil {
		data.UpdatedAt = t
	}
	return err
}

func (a AppsQuerier) Insert(data *App) error {
	data.UpdatedAt = time.Now()
	res, err := a.db.InsertInto(appTable).Columns(
		"type",
		"url",
		"name",
		"start_command",
		"kill_command",
		"is_system",
		"extra_kill",
		"announce",
		"require_credentials",
		"is_installed",
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
