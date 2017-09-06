// UDP meeting point server
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/anacrolix/utp"
	flags "github.com/jessevdk/go-flags"
	"github.com/mh-cbon/rendez-vous/browser"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/server"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/utils"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
)

type mainOpts struct {
	Version bool `long:"version" description:"Show version"`
}

//todo: rendez-vous server should check ttl registrations
//todo: rendez-vous server should implement a write token concept to register
//todo: rendez-vous server unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset}: %{message}`,
)

func showVersion() {
	fmt.Println("rendez-vous - noversion")
}

var options mainOpts
var parser = flags.NewParser(&options, flags.None)

func main() {
	parser.Parse()
	if options.Version {
		showVersion()
		os.Exit(0)
	}

	var cmds = flags.NewNamedParser("commands", flags.Default)
	cmds.AddCommand("serve",
		"Run rendez-vous server",
		"The serve command initialize a rendez-vous server which peers can use to register/unregister/find.",
		&rendezVousServerCommand{})

	cmds.AddCommand("client",
		"Run rendez-vous client",
		"The client command let you perform query on given remote rendez-vous server.",
		&rendezVousClientCommand{})

	cmds.AddCommand("website",
		"Run and announce a website on given rendez-vous remote.",
		"The website command runs a website and announce it to given remote rendez-vous server.",
		&rendezVousWebsiteCommand{})

	cmds.AddCommand("browser",
		"Run a browser to visit websites within a rendez-vous network.",
		"Starts a browser with a special local proxy that adequatly forwards incoming http requests on the network.",
		&rendezVousBrowserCommand{})

	cmds.AddCommand("http",
		"Run an http request using a rendez-vous client.",
		"Executes http requests over utp.",
		&rendezVousHTTPCommand{})

	if _, err := cmds.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}

type rendezVousServerCommand struct {
	Listen string `short:"l" long:"listen" description:"Port to listen" default:"0"`
}

func (opts *rendezVousServerCommand) Execute(args []string) error {
	if opts.Listen == "" {
		return fmt.Errorf("-listen argument is required")
	}

	conn, err := utils.UDP(":" + opts.Listen)
	if err != nil {
		return err
	}

	srv := socket.FromConn(conn)
	c := client.FromSocket(srv)
	srv.Handle(server.HandleQuery(c, nil))

	done := make(chan error)
	go handleSignal(done, srv.Close)

	readyErr := ready(func() error {
		log.Println("Listening...", ":"+opts.Listen)
		return nil
	}, srv.ListenAndServe)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

type rendezVousClientCommand struct {
	Listen string `short:"l" long:"listen" description:"Port to listen" default:"0"`
	Remote string `short:"r" long:"remote" description:"The rendez-vous address"`
	Query  string `short:"q" long:"query" description:"The query verb to run"`
	Pbk    string `long:"pbk" description:"An ed25519 prublic key - 32 len hex"`
	Pvk    string `long:"pvk" description:"The ed25519 private key - 64 len hex"`
	Value  string `long:"value" description:"The value to sign"`
	Retry  int    `long:"retry" description:"retry count"`
}

func (opts *rendezVousClientCommand) Execute(args []string) error {
	if opts.Listen == "" {
		return fmt.Errorf("-listen argument is required")
	}
	if opts.Remote == "" {
		return fmt.Errorf("-remote argument is required")
	}
	if model.OkVerb(opts.Query) == false {
		return fmt.Errorf("-query argument must be one of: %v", model.Verbs)
	}

	conn, err := utils.UDP(":" + opts.Listen)
	if err != nil {
		return err
	}

	srv := socket.FromConn(conn)
	c := client.FromSocket(srv)
	srv.Handle(client.HandleQuery(c))

	defer srv.Close()

	readyErr := ready(func() error {

		if opts.Query == "find" {
			id, err := identity.FromPbk(opts.Pbk, opts.Value)
			if err != nil {
				return errors.WithMessage(err, opts.Query)
			}
			res, err := c.Find(opts.Remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.Query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.Query == "unregister" {

			id, err := identity.FromPvk(opts.Pvk, opts.Value)
			if err != nil {
				return err
			}
			fmt.Println("pvk=", id.Pvk)
			fmt.Println("pbk=", id.Pbk)

			res, err := c.Unregister(opts.Remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.Query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.Query == "register" {

			id, err := identity.FromPvk(opts.Pvk, opts.Value)
			if err != nil {
				return err
			}
			fmt.Println("pvk=", id.Pvk)
			fmt.Println("pbk=", id.Pbk)
			fmt.Println("sig=", id.Sign)

			res, err := c.Register(opts.Remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.Query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.Query == "ping" {
			res, err := c.Ping(opts.Remote)
			if err != nil {
				if opts.Retry > 0 {
					for i := 0; i < opts.Retry; i++ {
						res, err = c.Ping(opts.Remote)
						if err == nil {
							break
						}
					}
				}
				if err != nil {
					return errors.WithMessage(err, opts.Query)
				}
			}
			fmt.Printf("%#v\n", res)

		} else {
			return errors.Errorf("Unknwon query %q", opts.Query)
		}
		return nil
	}, srv.ListenAndServe)

	return readyErr
}

type rendezVousWebsiteCommand struct {
	Listen string `short:"l" long:"listen" description:"Port to listen" default:"0"`
	Remote string `short:"r" long:"remote" description:"The rendez-vous address"`
	Local  string `long:"local" description:"The local port of the website" default:"9005"`
	Dir    string `long:"dir" description:"The directory of the me.com website" default:"demows"`
	Pvk    string `long:"pvk" description:"The ed25519 private key - 64 len hex - auto generated if empty"`
	Value  string `long:"value" description:"The value to sign" default:"website"`
}

func (opts *rendezVousWebsiteCommand) Execute(args []string) error {
	if opts.Listen == "" {
		return fmt.Errorf("-listen argument is required")
	}
	if opts.Dir == "" {
		return fmt.Errorf("-dir argument is required")
	}

	ln, err := utp.Listen(":" + opts.Listen)
	if err != nil {
		return err
	}

	id, err := identity.FromPvk(opts.Pvk, opts.Value)
	if err != nil {
		return err
	}

	pc := ln.(*utp.Socket)
	srv := socket.FromConn(pc)
	c := client.FromSocket(srv)
	srv.Handle(client.HandleQuery(c))

	registration := client.NewRegistration(time.Second*30, c)
	registration.Config(opts.Remote, *id)

	handler := http.FileServer(http.Dir(opts.Dir))
	public := utils.ServeHTTPFromListener(ln, httpServer(handler, "")) //todo: replace with a transparent proxy, so the website can live into another process
	local := httpServer(handler, "127.0.0.1:"+opts.Local)

	done := make(chan error)
	go handleSignal(done, registration.Stop, srv.Close, public.Close, local.Close)

	readyErr := ready(func() error {
		log.Println("Public Website listening on ", ln.Addr())
		log.Println("Local Website listening on ", local.Addr)

		fmt.Println("pvk=", id.Pvk)
		fmt.Println("pbk=", id.Pbk)
		fmt.Println("sig=", id.Sign)

		return nil
	}, srv.ListenAndServe, registration.Start, public.ListenAndServe, local.ListenAndServe)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

type rendezVousBrowserCommand struct {
	Listen   string `short:"l" long:"listen" description:"Port to listen" default:"0"`
	Remote   string `short:"r" long:"remote" description:"The rendez-vous address"`
	Proxy    string `short:"p" long:"proxy" description:"The port of the proxy" default:"9005"`
	Dir      string `long:"dir" description:"The directory of the website" default:"browser/static/"`
	Ws       string `short:"w" long:"ws" description:"The port of the website"`
	Headless bool   `long:"headless" description:"Run in headless mode (no-gui)"`
}

func (opts *rendezVousBrowserCommand) Execute(args []string) error {
	if opts.Listen == "" {
		return fmt.Errorf("--listen argument is required")
	}
	if opts.Remote == "" {
		return fmt.Errorf("--remote argument is required")
	}
	if opts.Proxy == "" {
		return fmt.Errorf("--proxy argument is required")
	}
	if opts.Ws == "" {
		return fmt.Errorf("--ws argument is required")
	}
	if opts.Dir == "" {
		return fmt.Errorf("--dir argument is required")
	}

	ln, err := utp.Listen(":" + opts.Listen)
	if err != nil {
		return err
	}

	pc := ln.(*utp.Socket)
	srv := socket.FromConn(pc)
	c := client.FromSocket(srv)

	wsAddr := "127.0.0.1:" + opts.Ws
	wsHandler := browser.MakeWebsite(opts.Dir)
	gateway := httpServer(wsHandler, wsAddr)

	browserProxy := browser.MakeProxyForBrowser(opts.Remote, wsAddr, c)
	proxy := httpServer(browserProxy, "127.0.0.1:"+opts.Proxy)

	done := make(chan error)
	go handleSignal(done, srv.Close, proxy.Close, gateway.Close)

	readyErr := ready(func() error {
		log.Println("me.com server listening on", wsAddr)

		if opts.Headless == false {
			cmd := exec.Command("chromium-browser", "--proxy-server="+proxy.Addr, "me.com")
			if err := cmd.Start(); err != nil {
				return err
			}
			cmd.Process.Release()
		}

		return nil
	}, srv.ListenAndServe, gateway.ListenAndServe, proxy.ListenAndServe)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

type rendezVousHTTPCommand struct {
	URL    string `short:"u" long:"url" description:"URL to execute"`
	Listen string `short:"l" long:"listen" description:"UTP port to listen" default:"0"`
	Remote string `short:"r" long:"remote" description:"The rendez-vous address"`
	Pbk    string `long:"pbk" description:"An ed25519 prublic key - 32 len hex"`
	Value  string `long:"value" description:"The value to sign" default:"website"`
	Knock  bool   `long:"knock" description:"Knock peer at first"`
}

func (opts *rendezVousHTTPCommand) Execute(args []string) error {
	if opts.URL == "" {
		return fmt.Errorf("-url argument is required")
	}
	if opts.Listen == "" {
		return fmt.Errorf("--listen argument is required")
	}
	if opts.Remote == "" {
		return fmt.Errorf("--remote argument is required")
	}

	ln, err := utp.Listen(":" + opts.Listen)
	if err != nil {
		return err
	}

	pc := ln.(*utp.Socket)
	srv := socket.FromConn(pc)
	c := client.FromSocket(srv)
	srv.Handle(client.HandleQuery(c))

	u, err := url.Parse(opts.URL)
	if err != nil {
		return errors.WithMessage(err, "url parse")
	}

	if opts.Knock {
		id, err2 := identity.FromPbk(opts.Pbk, opts.Value)
		if err2 != nil {
			return fmt.Errorf("knock failure: %v", err2.Error())
		}
		found, err2 := c.Knock(opts.Remote, id)
		if err2 != nil {
			for i := 0; i < 5; i++ {
				_, err2 = c.Ping(found.Data)
				if err2 == nil {
					break
				}
				<-time.After(time.Second)
			}
		}
		if err2 != nil {
			return fmt.Errorf("knock failure: %v", err2.Error())
		}
		u.Host = found.Data
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				dial, err2 := pc.Dial(u.Host)
				if err2 != nil {
					return nil, errors.WithMessage(err2, "dial")
				}
				return dial, nil
			},
		},
	}
	res, err := httpClient.Get(u.String())
	if err != nil {
		return errors.WithMessage(err, "http get")
	}
	defer res.Body.Close()
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil && err != io.EOF {
		return errors.WithMessage(err, "copy response")
	}
	return nil
}

func handleSignal(done chan error, do ...func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	var err error
	for _, d := range do {
		if e := d(); e != nil {
			err = e
		}
	}
	if done != nil {
		done <- err
	}
}

func httpServer(r http.Handler, addr string) *http.Server {
	return &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func ready(do func() error, blockings ...func() error) error {
	for index := range blockings {
		b := blockings[index]
		err := timeout(b, time.Millisecond*10)
		if err != nil {
			return err
		}
	}
	return do()
}

func timeout(do func() error, d time.Duration) error {
	rcv := make(chan error)
	go func() {
		rcv <- do()
	}()
	select {
	case err := <-rcv:
		close(rcv)
		if err != nil {
			return err
		}
	case <-time.After(d):
	}
	return nil
}
