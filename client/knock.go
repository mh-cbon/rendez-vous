package client

import (
	"errors"
	"sync"
	"time"
)

type Knock struct {
	id     string
	remote string
	create time.Time
	done   chan string
}

type ResolveHandler func(remote string)

func NewKnock(remote string, id string) Knock {
	if id == "" {
		id = "random"
	}
	return Knock{
		remote: remote,
		id:     id, //todo: random token
		create: time.Now(),
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

func (k Knock) Run(c *Client) (remote string, err error) {
	for i := 0; i < 5; i++ {
		go c.Knock(k.remote, k.id)
		select {
		case res := <-k.done:
			return res, nil
		case <-time.After(time.Second):
		}
	}
	return "", errors.New("knock timeout")
}

type PendingKnocks struct {
	knocks map[string]Knock
}

func NewPendingKnocks() *PendingKnocks {
	return &PendingKnocks{
		knocks: map[string]Knock{},
	}
}

func (p *PendingKnocks) Add(remote string, id string) Knock {
	k := NewKnock(remote, id)
	p.knocks[k.id] = k
	return k
}
func (p *PendingKnocks) Rm(k Knock) bool {
	if _, ok := p.knocks[k.id]; ok {
		delete(p.knocks, k.id)
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

type PendingKnocksTS struct {
	store *PendingKnocks
	l     sync.Mutex
}

func NewPendingKnocksTS(store *PendingKnocks) *PendingKnocksTS {
	if store == nil {
		store = NewPendingKnocks()
	}
	return &PendingKnocksTS{
		store: store,
		l:     sync.Mutex{},
	}
}

func (p *PendingKnocksTS) Add(remote string, id string) Knock {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Add(remote, id)
}

func (p *PendingKnocksTS) Rm(k Knock) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Rm(k)
}

func (p *PendingKnocksTS) Resolve(remote string, id string) bool {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Resolve(remote, id)
}
