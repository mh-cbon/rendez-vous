package shellexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Command Return a new exec.Cmd object for the given command string
func Command(cwd string, cmd string) (*Cmd, error) {
	return NewTempCmd(cwd, cmd)
}

// Cmd ...
type Cmd struct {
	*exec.Cmd
	f string
}

// NewTempCmd ...
func NewTempCmd(cwd string, cmd string) (*Cmd, error) {
	f, err := ioutil.TempDir("", "shellexec")
	if err != nil {
		return nil, err
	}
	fp := filepath.Join(f, "s")
	if isWindows {
		fp += ".bat"
	}
	err = ioutil.WriteFile(fp, []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}
	ret := &Cmd{Cmd: exec.Command("sh", "-c", fp), f: fp}
	if isWindows {
		ret.Cmd = exec.Command("cmd", "/C", fp)
	}
	ret.Cmd.Dir = cwd
	return ret, nil
}

// Start ...
func (t *Cmd) Start() error {
	return t.Cmd.Start()
}

// Run ...
func (t *Cmd) Run() error {
	if err := t.Start(); err != nil {
		return err
	}
	return t.Wait()
}

// Wait ...
func (t *Cmd) Wait() error {
	err := t.Cmd.Wait()
	os.Remove(t.f)
	return err
}
