package admin

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/mh-cbon/rendez-vous/admin/shellexec"
)

type IOProvider struct {
	Stdout func(name, cwd, cmdstr string) io.Writer
	Stderr func(name, cwd, cmdstr string) io.Writer
}

var StdIO = &IOProvider{
	func(name, cwd, cmdstr string) io.Writer { return os.Stdout },
	func(name, cwd, cmdstr string) io.Writer { return os.Stderr },
}

var BufferedIO = &IOProvider{
	func(name, cwd, cmdstr string) io.Writer { return new(bytes.Buffer) },
	func(name, cwd, cmdstr string) io.Writer { return new(bytes.Buffer) },
}

type ChildProcesses struct {
	processes map[string]*shellexec.Cmd
	io        *IOProvider
}

func NewChildProcesses(ioProvider *IOProvider) *ChildProcesses {
	if ioProvider == nil {
		ioProvider = StdIO
	}
	return &ChildProcesses{
		processes: map[string]*shellexec.Cmd{},
		io:        ioProvider,
	}
}

func (c *ChildProcesses) Get(name string) *shellexec.Cmd {
	if x, ok := c.processes[name]; ok {
		return x
	}
	return nil
}

func (c *ChildProcesses) Start(name, cwd, cmdstr string, env ...string) error {
	if c.Get(name) != nil {
		return fmt.Errorf("%q already started", name)
	}
	cp, err := shellexec.Command(cwd, cmdstr)
	if err != nil {
		return err
	}
	cp.Stdout = c.io.Stdout(name, cwd, cmdstr)
	cp.Stderr = c.io.Stderr(name, cwd, cmdstr)
	cp.Env = append(cp.Env, env...)
	if err = cp.Start(); err != nil {
		return err
	}
	c.processes[name] = cp
	return nil
}

func (c *ChildProcesses) Run(name, cwd, cmdstr string, env ...string) error {
	if err := c.Start(name, cwd, cmdstr, env...); err != nil {
		return err
	}
	return c.Wait(name)
}

func (c *ChildProcesses) Wait(name string) error {
	cp := c.Get(name)
	if cp == nil {
		return fmt.Errorf("%q not started", name)
	}
	if err := cp.Wait(); err != nil {
		return err
	}
	delete(c.processes, name)
	return nil
}

func (c *ChildProcesses) Kill(name string) error {
	cp := c.Get(name)
	if cp == nil {
		return fmt.Errorf("%q not started", name)
	}
	if err := cp.Process.Kill(); err != nil {
		return err
	}
	delete(c.processes, name)
	return nil
}

func (c *ChildProcesses) Pid(name string) (int, error) {
	cp := c.Get(name)
	if cp != nil {
		return 0, fmt.Errorf("%q not started", name)
	}
	return cp.Process.Pid, nil
}
