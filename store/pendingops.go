package store

import (
	"errors"
	"sync"
	"time"

	"github.com/mh-cbon/rendez-vous/model"
)

type PendingOp struct {
	Token             string
	Create            time.Time
	resolutionHandler ResolutionHandler
	done              chan Resolution
}

type ResolutionHandler func(remote string, m model.Message) bool

type Resolution struct {
	Remote string
	M      model.Message
}

func NewPendingOp(token string, resolutionHandler ResolutionHandler) *PendingOp {
	if token == "" {
		token = "random" //todo: random token
	}
	return &PendingOp{
		Token:             token,
		Create:            time.Now(),
		resolutionHandler: resolutionHandler,
		done:              make(chan Resolution),
	}
}

func (k *PendingOp) Resolve(remote string, m model.Message) bool {
	if k.resolutionHandler(remote, m) {
		go func() {
			k.done <- Resolution{remote, m}
		}()
		return true
	}
	return false
}

func (k *PendingOp) Done() <-chan Resolution {
	return k.done
}

func (k *PendingOp) Run(h func() error) (Resolution, error) {
	w := make(chan error)
	var err error
	var res Resolution
	go func() {
		w <- h()
	}()
	select {
	case res = <-k.Done():
	case err = <-w:
	case <-time.After(time.Second):
		err = errors.New("op timeout")
	}
	return res, err
}

type PendingOps struct {
	ops map[string]*PendingOp
}

func NewPendingOps() *PendingOps {
	return &PendingOps{
		ops: map[string]*PendingOp{},
	}
}

func (p *PendingOps) Add(token string, resolutionHandler ResolutionHandler) *PendingOp {
	k := NewPendingOp(token, resolutionHandler)
	p.ops[token] = k
	return k
}
func (p *PendingOps) Rm(k *PendingOp) bool {
	token := k.Token
	if _, ok := p.ops[token]; ok {
		delete(p.ops, token)
		return true
	}
	return false
}
func (p *PendingOps) Resolve(remote string, token string, m model.Message) bool {
	if x, ok := p.ops[token]; ok && x.Resolve(remote, m) {
		delete(p.ops, token)
		return true
	}
	return false
}

type TSPendingOps struct {
	store *PendingOps
	l     sync.Mutex
}

func NewTSPendingOps(store *PendingOps) *TSPendingOps {
	if store == nil {
		store = NewPendingOps()
	}
	return &TSPendingOps{
		store: store,
		l:     sync.Mutex{},
	}
}

func (p *TSPendingOps) Add(token string, resolutionHandler ResolutionHandler) *PendingOp {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Add(token, resolutionHandler)
}

func (p *TSPendingOps) Rm(k *PendingOp) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Rm(k)
}

func (p *TSPendingOps) Resolve(remote string, token string, m model.Message) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Resolve(remote, token, m)
}
