package main

import (
	"flag"
	"net/http"
)

func main() {
	port := flag.String("port", "8084", "")
	dir := flag.String("dir", "demows", "")
	flag.Parse()
	p := *port
	d := *dir
	http.ListenAndServe(":"+p, http.FileServer(http.Dir(d)))
}
