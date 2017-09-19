package main

import (
	"flag"
	"net/http"

	"github.com/mh-cbon/rendez-vous/cli"
)

func main() {

	cli.RendezVous(cli.App{
		Start: "go run main.go -listen {{.listen}} -dir {{.dir}}",
		StartParams: []cli.Param{
			cli.Param{Name: "listen", Default: ":8090", Type: "netaddr"},
			cli.Param{Name: "dir", Default: "assets", Type: "dirpath"},
		},
		// Kill:        "",
		// KillParams:  []cli.Param{},
	})

	var listen string
	var dir string

	flag.StringVar(&listen, ":8090", "listen", "Listen address")
	flag.StringVar(&dir, "assets", "dir", "Directory to serve")
	flag.Parse()

	handler := http.FileServer(http.Dir(dir))
	err := http.ListenAndServe(listen, handler)
	if err != nil {
		panic(err)
	}
}
