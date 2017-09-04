package browser

import (
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/anacrolix/utp"
	"github.com/elazarl/goproxy"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/identity"
)

// MakeProxyForBrowser prepares a proxy to handle me.com requests
// if the request is in the form me.com/...
// it forwards the query to the given wsAddr
// if the request is in the form <pbk>.me.com/...
// then it searches for the peer address on remote
// if found, it proxy the http request to the peer found.
func MakeProxyForBrowser(remote string, wsAddr string, c *client.Client) http.Handler {
	browserProxy := goproxy.NewProxyHttpServer()
	browserProxy.Verbose = true
	browserProxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.+[.]me[.]com$"))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			pbk := strings.Split(r.URL.Host, ".")[0]

			id, err := identity.FromPbk(pbk, "website")
			if err != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"failed: "+err.Error())
			}

			peer, err := c.Find(remote, id)
			if err != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"failed: "+err.Error())
			}
			log.Println(peer)

			httpClient := http.Client{
				Transport: &http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						return utp.Dial(peer.Response)
					},
				},
			}

			r.URL.Host = peer.Response
			res, err := httpClient.Get(r.URL.String())
			if err != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"failed: "+err.Error())
			}
			return nil, res
		})
	browserProxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^me[.]com$"))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.URL.Host = wsAddr
			return r, nil
		})

	return browserProxy
}