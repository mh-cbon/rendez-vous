// UDP meeting point server
package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/anacrolix/utp"
	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/client"
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
	value  string
	sign   string
	auto   bool
}

type websiteOpts struct {
	listen string
	remote string
	dir    string
	pbk    string
	value  string
	sign   string
	auto   bool
}

type browserOpts struct {
	listen string
	remote string
	proxy  string
	ws     string
}

//todo: add storage clean up with ttl on entry

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset}: %{message}`,
)

func showErr(flags *flag.FlagSet, reason string) {
	showHelp(flags)
	fmt.Println()
	fmt.Println("Wrong command line: ", reason)
}
func showHelp(flags *flag.FlagSet) {
	showVersion()
	fmt.Println("")
	fmt.Println("	A server to expose your endpoints with a public key.")
	fmt.Println("	A client to find/register endpoints for a public key.")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("	rendez-vous [server|client] <options>")
	if flags != nil {
		fmt.Println("")
		flags.Usage()
	}
}
func showVersion() {
	fmt.Println("rendez-vous - noversion")
}

func main() {

	sHelp := flag.Bool("h", false, "show help")
	lHelp := flag.Bool("help", false, "show help")

	sVersion := flag.Bool("v", false, "show version")
	lVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	if *sHelp || *lHelp {
		return
	} else if *sVersion || *lVersion {
		showVersion()
		return
	} else if len(os.Args) < 2 {
		showErr(flag.CommandLine, "Missing operation (server|client)")
		return
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

		runServer(opts)

	} else if op == "client" {

		var opts cliOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Port to listen")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.StringVar(&opts.query, "query", "", "Query to send (ping|find|register|unregister)")
		set.BoolVar(&opts.auto, "auto", false, "Generate pvk/pbk/sign automatically")
		set.StringVar(&opts.pbk, "pbk", "", "Pbk of the query in hexadecimal")
		set.StringVar(&opts.value, "value", "", "Value of the query")
		set.StringVar(&opts.sign, "sign", "", "Sign of the query in hexadecimal")
		sHelp = set.Bool("-h", false, "show help")
		lHelp = set.Bool("-help", false, "show help")
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
		runClient(opts)

	} else if op == "website" {
		var opts websiteOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Port to listen")
		set.StringVar(&opts.dir, "static", "./static", "Directory to serve")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.BoolVar(&opts.auto, "auto", false, "Generate pvk/pbk/sign automatically")
		set.StringVar(&opts.pbk, "pbk", "", "Pbk of the query in hexadecimal")
		set.StringVar(&opts.value, "value", "", "Value of the query")
		set.StringVar(&opts.sign, "sign", "", "Sign of the query in hexadecimal")
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

		runWebsite(opts)

	} else if op == "browser" {
		var opts browserOpts
		set := flag.NewFlagSet(op, flag.ExitOnError)
		set.StringVar(&opts.listen, "listen", "0", "Port to listen")
		set.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
		set.StringVar(&opts.proxy, "proxy", "", "Address of the proxy")
		set.StringVar(&opts.ws, "ws", "", "Address of the local website")
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

		runBrowser(opts)

	} else {
		log.Fatal("Wrong command line, must be: rendez-vous [server|client|website|browser] <options>")
	}
}

func runServer(opts srvOpts) {

	conn, err := utils.UDP(":" + opts.listen)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.FromSocket(socket.FromConn(conn))
	if err := srv.Listen(); err != io.EOF {
		log.Fatal(err)
	}
}

func runClient(opts cliOpts) {

	conn, err := utils.UDP(":" + opts.listen)
	if err != nil {
		log.Fatal(err)
	}

	c := client.FromSocket(socket.FromConn(conn))

	{
		var res model.Message
		var err error
		if opts.query == "find" {
			res, err = c.Find(opts.remote, opts.pbk)

		} else if opts.query == "unregister" {
			res, err = c.Unregister(opts.remote, opts.pbk)

		} else if opts.query == "register" {
			if opts.auto {
				pvk, _, err2 := ed25519.GenerateKey(rand.Reader)
				if err2 != nil {
					log.Fatal(err2)
				}
				pbk := ed25519.PublicKeyFromPvk(pvk)
				sign := ed25519.Sign(pvk, pbk, []byte(opts.value))
				opts.pbk = hex.EncodeToString(pbk)
				opts.sign = hex.EncodeToString(sign)
			}
			res, err = c.Register(opts.remote, opts.pbk, opts.sign, opts.value)

		} else if opts.query == "ping" {
			res, err = c.Ping(opts.remote)
			if err != nil {
				err = errors.WithMessage(err, "query ping")
			}

		} else {
			err = fmt.Errorf("Unknwon query %q", opts.query)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", res)
	}
}

func runWebsite(opts websiteOpts) {

	if opts.auto {
		pvk, _, err2 := ed25519.GenerateKey(rand.Reader)
		if err2 != nil {
			log.Fatal(err2)
		}
		pbk := ed25519.PublicKeyFromPvk(pvk)
		sign := ed25519.Sign(pvk, pbk, []byte(opts.value))
		opts.pbk = hex.EncodeToString(pbk)
		opts.sign = hex.EncodeToString(sign)
	}

	ln, err := utp.Listen(":" + opts.listen)
	if err != nil {
		log.Fatal(err)
	}

	pc := ln.(*utp.Socket)
	c := client.FromSocket(socket.FromConn(pc))
	_, err = c.Register(opts.remote, opts.pbk, opts.pbk, opts.value)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: http.FileServer(http.Dir(opts.dir)),
	}
	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}

func runBrowser(opts browserOpts) {

	ln, err := utp.Listen(":" + opts.listen)
	if err != nil {
		log.Fatal(err)
	}

	pc := ln.(*utp.Socket)
	c := client.FromSocket(socket.FromConn(pc))

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./browser/static/")))
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + opts.ws,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening...", "127.0.0.1:"+opts.ws)

	go func() {
		if err2 := srv.ListenAndServe(); err2 != nil {
			log.Fatal(err2)
		}
	}()

	go func() {
		proxy := goproxy.NewProxyHttpServer()
		proxy.Verbose = true
		proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.+[.]me[.]com$"))).DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				pbk := strings.Split(r.URL.Host, ".")[0]
				res, err := c.Find(opts.remote, pbk)
				if err != nil {
					return r, goproxy.NewResponse(r,
						goproxy.ContentTypeText, http.StatusForbidden,
						"failed: "+err.Error())
				}
				log.Println(res)
				return r, nil
			})
		proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^me[.]com$"))).DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				r.URL.Host = "127.0.0.1:" + opts.ws
				return r, nil
			})
		log.Fatal(http.ListenAndServe("127.0.0.1:"+opts.proxy, proxy))
	}()

	go func() {
		cmd := exec.Command("chromium-browser", "--proxy-server=127.0.0.1:"+opts.proxy, "me.com")
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		cmd.Process.Release()
	}()

	<-make(chan bool)
}
