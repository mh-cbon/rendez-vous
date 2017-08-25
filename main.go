// UDP meeting point server
package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mh-cbon/dht/ed25519"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/server"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
	"github.com/pkg/errors"
)

type cliOpts struct {
	op     string
	port   string
	remote string
	query  string
	pbk    string
	value  string
	sign   string
	auto   bool
}

//todo: add storage clean up with ttl on entry

func main() {

	var opts cliOpts

	flag.StringVar(&opts.op, "op", "", "operation to do server|client")
	flag.StringVar(&opts.port, "port", "8080", "Port to listen")
	flag.StringVar(&opts.remote, "remote", "", "Address of the rendez-vous")
	flag.StringVar(&opts.query, "query", "", "Query to send")
	flag.BoolVar(&opts.auto, "auto", false, "Generate pvk/pbk/sign automatically")
	flag.StringVar(&opts.pbk, "pbk", "", "Pbk of the query in hexadecimal")
	flag.StringVar(&opts.value, "value", "", "Value of the query")
	flag.StringVar(&opts.sign, "sign", "", "Sign of the query in hexadecimal")

	flag.Parse()

	switch opts.op {
	case "server":
		runServer(opts)
	case "client":
		runClient(opts)
	default:
		log.Fatal("Wrong command line, must be: rendez-vous [server|client] ...options")
	}
}

func runServer(opts cliOpts) {

	if opts.port == "" {
		log.Fatalf("-port argument is required")
	}

	storage := store.New(nil)

	s, err := socket.FromAddr(":" + opts.port)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Listen(server.Handler(storage)); err != io.EOF {
		log.Fatal(err)
	}
}

func runClient(opts cliOpts) {

	if opts.remote == "" {
		log.Fatalf("-remote argument is required")
	}

	if opts.port == "" {
		log.Fatalf("-remote argument is required")
	}

	remote, err := net.ResolveUDPAddr("udp", opts.remote)
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.FromAddr(":" + opts.port)
	if err != nil {
		log.Fatal(err)
	}

	persist := false
	{
		var res model.Message
		var err error
		if opts.query == "find" {
			b, err2 := hex.DecodeString(opts.pbk)
			if err2 != nil {
				log.Fatal(err2)
			}
			res, err = c.Find(remote, b)

		} else if opts.query == "unregister" {
			b, err2 := hex.DecodeString(opts.pbk)
			if err2 != nil {
				log.Fatal(err2)
			}
			res, err = c.Unregister(remote, b)

		} else if opts.query == "register" {
			var pbk []byte
			var sign []byte
			if opts.auto {
				pvk, _, err2 := ed25519.GenerateKey(rand.Reader)
				if err2 != nil {
					log.Fatal(err2)
				}

				pbk = ed25519.PublicKeyFromPvk(pvk)
				sign = ed25519.Sign(pvk, pbk, []byte(opts.value))

				fmt.Printf("pvk %x\n", pvk)
				fmt.Printf("pbk %x\n", pbk)
				fmt.Printf("sign %x\n", sign)
				fmt.Printf("sign %v\n", len(sign))

			} else {
				ppbk, err2 := hex.DecodeString(opts.pbk)
				if err2 != nil {
					log.Fatal(err2)
				}
				psign, err2 := hex.DecodeString(opts.sign)
				if err2 != nil {
					log.Fatal(err2)
				}
				pbk = ppbk
				sign = psign
			}
			res, err = c.Register(remote, pbk, sign, opts.value)

			persist = err == nil

		} else if opts.query == "ping" {
			res, err = c.Ping(remote)
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

		// only for demo
		if persist {
			var b [0x10000]byte
			for {
				n, _ := c.Conn().Read(b[:])
				if len(b) > 0 {
					fmt.Println(string(b[:n]))
				}
			}
		}
	}
}
