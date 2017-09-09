package browser

import (
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/node"
)

func NewProxy(nodeListen, remote, websiteListen, proxyListen string, me *identity.PublicIdentity) *Proxy {
	return &Proxy{
		node:          node.NewPeerNode(nodeListen),
		remote:        remote,
		websiteListen: websiteListen,
		proxyListen:   proxyListen,
		me:            me,
	}
}

type Proxy struct {
	node          *node.PeerNode
	remote        string
	websiteListen string
	proxy         *http.Server
	proxyListen   string
	me            *identity.PublicIdentity
}

func (r *Proxy) ListenAndServe() error {

	if err := r.node.Start(); err != nil {
		return err
	}

	node := r.node
	remote := r.remote
	me := r.me
	websiteListen := r.websiteListen

	browserProxy := goproxy.NewProxyHttpServer()
	browserProxy.Verbose = true
	browserProxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.+[.]me[.]com$"))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			httpClient := http.Client{
				Transport: &http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						host, err := node.Resolve(remote, addr, "website", me)
						if err != nil {
							return nil, err
						}
						return node.Dial(host, "website")
					},
				},
			}
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
			r.URL.Host = websiteListen
			return r, nil
		})
	r.proxy = httpServer(browserProxy, r.proxyListen)
	return r.proxy.ListenAndServe()
}

func (r *Proxy) Close() error {
	var err error
	if err2 := r.node.Close(); err2 != nil {
		err = err2
	}
	if err2 := r.proxy.Close(); err2 != nil {
		err = err2
	}
	return err
}

func (r *Proxy) ChangeListenAddress(listen string) error {
	return r.node.Restart(listen)
}

func (r *Proxy) LocalAddr() net.Addr {
	return r.node.LocalAddr()
}

func (r *Proxy) Port() *node.PortStatus {
	return r.node.Port
}

func (r *Proxy) Client() *client.Client {
	return r.node.GetClient()
}

func (r *Proxy) TestPort() *node.PortStatus {
	return r.node.TestPort(r.remote, nil)
}

func (r *Proxy) List(start, limit int) ([]*model.Peer, error) {
	return r.node.GetClient().List(r.remote, start, limit)
}

func httpServer(r http.Handler, addr string) *http.Server {
	return &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
