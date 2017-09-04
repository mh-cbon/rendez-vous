// UDP meeting point server
package main

import (
	"flag"
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

type srvOpts struct {
	listen string
}

type cliOpts struct {
	listen string
	remote string
	query  string
	pbk    string
	pvk    string
	value  string
	sign   string
	auto   bool
}

type websiteOpts struct {
	listen string
	local  string
	remote string
	dir    string
	pvk    string
	value  string
}

type browserOpts struct {
	listen   string
	remote   string
	proxy    string
	ws       string
	dir      string
	headless bool
}

type httpOpts struct {
	url    string
	method string
}

//todo: add storage clean up with ttl on entry

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset}: %{message}`,
)

func showErr(flags *flag.FlagSet, reason ...interface{}) {
	showHelp(flags)
	fmt.Println()
	fmt.Print("Wrong command line: ")
	fmt.Println(reason...)
}
func showHelp(flags *flag.FlagSet) {
	showVersion()
	fmt.Println("")
	fmt.Println("	A server to expose your endpoints with a public key.")
	fmt.Println("	A client to find/register endpoints for a public key.")
	fmt.Println("	A website to expose website.")
	fmt.Println("	A browser to visit remote website.")
	fmt.Println("	An http client to test a remote website..")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("	rendez-vous [server|client|website|browser|http] <options>")
	if flags != nil {
		fmt.Println("")
		flags.Usage()
	}
}
func showVersion() {
	fmt.Println("rendez-vous - noversion")
}

func main() {

	flag.CommandLine = flag.NewFlagSet("main", flag.ExitOnError) // get ride of all test* flags.
	sHelp := flag.Bool("h", false, "show help")
	lHelp := flag.Bool("help", false, "show help")

	sVersion := flag.Bool("v", false, "show version")
	lVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	if len(os.Args) < 2 {
		showHelp(flag.CommandLine)
		return

	} else if len(os.Args) < 3 {

		if *sHelp || *lHelp {
			showHelp(flag.CommandLine)
			return
		} else if *sVersion || *lVersion {
			showVersion()
			return
		}
	}

	op := os.Args[1]
	args := os.Args[2:]

	if op == "serve" {
		var opts srvOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Port to listen")
		sHelp = set.Bool("h", false, "show help")
		lHelp = set.Bool("help", false, "show help")
		set.Parse(args)
		if *sHelp || *lHelp {
			showHelp(set)
			return
		}
		if opts.listen == "" {
			showErr(set, "-listen argument is required")
			return
		}

		if err := runServer(opts); err != nil {
			log.Fatal(err)
		}

	} else if op == "client" {

		var opts cliOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Port to listen")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.StringVar(&opts.query, "query", "", "Query to send (ping|find|register|unregister)")
		set.StringVar(&opts.pbk, "pbk", "", "Pbk to lookup for")
		set.StringVar(&opts.pvk, "pvk", "", "Pvk of the registration in hexadecimal")
		set.StringVar(&opts.value, "value", "", "Value of the query")
		sHelp = set.Bool("h", false, "show help")
		lHelp = set.Bool("help", false, "show help")
		set.Parse(args)
		if *sHelp || *lHelp {
			showHelp(set)
			return
		}
		if opts.listen == "" {
			showErr(set, "-listen argument is required")
			return
		}
		if opts.remote == "" {
			showErr(set, "-remote argument is required")
			return
		}
		if model.OkVerb(opts.query) == false {
			showErr(set, "-query argument is invalid")
			return
		}

		if err := runClient(opts); err != nil {
			log.Fatal(err)
		}

	} else if op == "website" {
		var opts websiteOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Public port")
		set.StringVar(&opts.local, "local", "9005", "Local port")
		set.StringVar(&opts.dir, "static", "./static", "Directory to serve")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.StringVar(&opts.pvk, "pvk", "", "Pvk used for registration, it is random if not empty")
		set.StringVar(&opts.value, "value", "website", "Value to to sign")
		sHelp = set.Bool("h", false, "show help")
		lHelp = set.Bool("help", false, "show help")
		set.Parse(args)
		if *sHelp || *lHelp {
			showHelp(set)
			return
		}
		if opts.listen == "" {
			showErr(set, "-listen argument is required")
			return
		}
		if opts.dir == "" {
			showErr(set, "-static argument is required")
			return
		}

		if err := runWebsite(opts); err != nil {
			log.Fatal(err)
		}

	} else if op == "browser" {
		var opts browserOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Public port to listen")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.StringVar(&opts.proxy, "proxy", "", "Address of the proxy")
		set.StringVar(&opts.ws, "ws", "", "Address of the local me.com website")
		set.StringVar(&opts.dir, "dir", "browser/static/", "Directory of the static assets for me.com")
		set.BoolVar(&opts.headless, "headless", false, "Headless mode (no gui)")
		sHelp = set.Bool("h", false, "show help")
		lHelp = set.Bool("help", false, "show help")
		set.Parse(args)
		if *sHelp || *lHelp {
			showHelp(set)
			return
		}
		if opts.listen == "" {
			showErr(set, "-listen argument is required")
			return
		}
		if opts.remote == "" {
			showErr(set, "-remote argument is required")
			return
		}
		if opts.proxy == "" {
			showErr(set, "-proxy argument is required")
			return
		}
		if opts.ws == "" {
			showErr(set, "-ws argument is required")
			return
		}

		if err := runBrowser(opts); err != nil {
			log.Fatal(err)
		}

	} else if op == "http" {
		var opts httpOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.url, "url", "", "URL to fetch http://ip:port/path")
		sHelp = set.Bool("h", false, "show help")
		lHelp = set.Bool("help", false, "show help")
		set.Parse(args)
		if *sHelp || *lHelp {
			showHelp(set)
			return
		}
		if opts.url == "" {
			showErr(set, "-url argument is required")
			return
		}

		if err := runHTTPClient(opts); err != nil {
			log.Fatal(err)
		}

	} else {
		showErr(flag.CommandLine, "Wrong werb ", op)
	}
}

func runServer(opts srvOpts) error {

	conn, err := utils.UDP(":" + opts.listen)
	if err != nil {
		return err
	}

	done := make(chan error)
	go handleSignal(done, conn.Close)

	srv := server.FromSocket(socket.FromConn(conn))
	readyErr := ready(func() error {
		log.Println("Listening...", ":"+opts.listen)
		return nil
	}, srv)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

func runClient(opts cliOpts) error {

	conn, err := utils.UDP(":" + opts.listen)
	if err != nil {
		return err
	}

	go handleSignal(nil, conn.Close)

	c := client.FromSocket(socket.FromConn(conn))

	readyErr := ready(func() error {

		if opts.query == "find" {
			id, err := identity.FromPbk(opts.pbk, opts.value)
			if err != nil {
				return errors.WithMessage(err, opts.query)
			}
			res, err := c.Find(opts.remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.query == "unregister" {

			id, err := identity.FromPvk(opts.pvk, opts.value)
			if err != nil {
				return err
			}
			fmt.Println("pvk=", id.Pvk)
			fmt.Println("pbk=", id.Pbk)

			res, err := c.Unregister(opts.remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.query == "register" {

			id, err := identity.FromPvk(opts.pvk, opts.value)
			if err != nil {
				return err
			}
			fmt.Println("pvk=", id.Pvk)
			fmt.Println("pbk=", id.Pbk)
			fmt.Println("sig=", id.Sign)

			res, err := c.Register(opts.remote, id)
			if err != nil {
				return errors.WithMessage(err, opts.query)
			}
			fmt.Printf("%#v\n", res)

		} else if opts.query == "ping" {
			res, err := c.Ping(opts.remote)
			if err != nil {
				return errors.WithMessage(err, opts.query)
			}
			fmt.Printf("%#v\n", res)

		} else {
			return errors.Errorf("Unknwon query %q", opts.query)
		}
		return nil
	}, c)
	if readyErr != nil {
		return readyErr
	}
	return nil
}

func runWebsite(opts websiteOpts) error {

	ln, err := utp.Listen(":" + opts.listen)
	if err != nil {
		return err
	}

	pc := ln.(*utp.Socket)
	c := client.FromSocket(socket.FromConn(pc))

	done := make(chan error)
	go handleSignal(done, ln.Close)

	handler := http.FileServer(http.Dir(opts.dir))
	public := httpu{httpServer(handler, ""), ln}
	local := httpServer(handler, "127.0.0.1:"+opts.local)

	readyErr := ready(func() error {
		log.Println("Public Website listening on ", ln.Addr())
		log.Println("Local Website listening on ", local.Addr)

		id, err := identity.FromPvk(opts.pvk, opts.value)
		if err != nil {
			return err
		}
		fmt.Println("pvk=", id.Pvk)
		fmt.Println("pbk=", id.Pbk)
		fmt.Println("sig=", id.Sign)

		res, err := c.Register(opts.remote, id)
		if err != nil {
			return err
		}

		log.Println("registration ", res)
		return err
	}, c, public, local)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

func runBrowser(opts browserOpts) error {

	ln, err := utp.Listen(":" + opts.listen)
	if err != nil {
		return err
	}

	pc := ln.(*utp.Socket)
	c := client.FromSocket(socket.FromConn(pc))

	done := make(chan error)
	go handleSignal(done, ln.Close)

	registration := client.NewRegistration(time.Second*5, c)

	wsAddr := "127.0.0.1:" + opts.ws
	wsHandler := browser.MakeWebsite(opts.dir)
	gateway := httpServer(wsHandler, wsAddr)

	browserProxy := browser.MakeProxyForBrowser(opts.remote, wsAddr, c)
	proxy := httpServer(browserProxy, "127.0.0.1:"+opts.proxy)

	readyErr := ready(func() error {
		log.Println("me.com server listening on", "127.0.0.1:"+opts.ws)

		if opts.headless == false {
			cmd := exec.Command("chromium-browser", "--proxy-server="+proxy.Addr, "me.com")
			if err := cmd.Start(); err != nil {
				return err
			}
			cmd.Process.Release()
		}

		return nil
	}, c, registration, gateway, proxy)
	if readyErr != nil {
		return readyErr
	}
	return <-done
}

func runHTTPClient(opts httpOpts) error {

	u, err := url.Parse(opts.url)
	if err != nil {
		return errors.WithMessage(err, "url parse")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return utp.Dial(u.Host)
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

type httpu struct {
	*http.Server
	l net.Listener
}

func (h httpu) ListenAndServe() error {
	return h.Server.Serve(h.l)
}

type listenAndServe interface {
	ListenAndServe() error
}

type listen interface {
	Listen() error
}

type start interface {
	Start()
}

func ready(do func() error, blockings ...interface{}) error {
	for index := range blockings {
		b := blockings[index]
		err := timeout(func() error {
			if x, ok := b.(listen); ok {
				return x.Listen()
			} else if x, ok := b.(listenAndServe); ok {
				return x.ListenAndServe()
			} else if x, ok := b.(start); ok {
				x.Start()
				return nil
			}
			return fmt.Errorf("unknown blocking interface for %T %#v", b, b)
		}, time.Millisecond*10)
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
