package store

import (
	"errors"
	"sync"
	"time"
)

type Knock struct {
	ID     string
	Create time.Time
	done   chan string
}

type ResolveHandler func(remote string)

func NewKnock(id string) Knock {
	if id == "" {
		id = "random"
	}
	return Knock{
		ID:     id, //todo: random token
		Create: time.Now(),
		done:   make(chan string),
	}
}

func (k Knock) Resolve(remote string) bool {
	// remote stringS must match
	go func() {
		k.done <- remote
	}()
	return true
}

func (k Knock) Done() <-chan string {
	return k.done
}

func (k Knock) Run(h func() error) (string, error) {
	w := make(chan error)
	var err error
	var res string
	go func() {
		w <- h()
	}()
	select {
	case res = <-k.done:
	case err = <-w:
	case <-time.After(time.Second):
		err = errors.New("knock timeout")
	}
	return res, err
}

type PendingKnocks struct {
	knocks map[string]Knock
}

func NewPendingKnocks() *PendingKnocks {
	return &PendingKnocks{
		knocks: map[string]Knock{},
	}
}

func (p *PendingKnocks) Add(id string) Knock {
	k := NewKnock(id)
	p.knocks[k.ID] = k
	return k
}
func (p *PendingKnocks) Rm(k Knock) bool {
	if _, ok := p.knocks[k.ID]; ok {
		delete(p.knocks, k.ID)
		return true
	}
	return false
}
func (p *PendingKnocks) Resolve(remote string, id string) bool {
	if x, ok := p.knocks[id]; ok && x.Resolve(remote) {
		delete(p.knocks, id)
		return true
	}
	return false
}

type TSPendingKnocks struct {
	store *PendingKnocks
	l     sync.Mutex
}

func NewTSPendingKnocks(store *PendingKnocks) *TSPendingKnocks {
	if store == nil {
		store = NewPendingKnocks()
	}
	return &TSPendingKnocks{
		store: store,
		l:     sync.Mutex{},
	}
}

func (p *TSPendingKnocks) Add(id string) Knock {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Add(id)
}

func (p *TSPendingKnocks) Rm(k Knock) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Rm(k)
}

func (p *TSPendingKnocks) Resolve(remote string, id string) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Resolve(remote, id)
}
