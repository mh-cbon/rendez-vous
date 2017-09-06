package client

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Knock struct {
	id     string
	create time.Time
	done   chan string
}

type ResolveHandler func(remote string)

func NewKnock(id string) Knock {
	if id == "" {
		id = "random"
	}
	return Knock{
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

func (k Knock) Run(remote string, c *Client) (string, error) {
	x := make(chan error)
	f := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			// var wg sync.WaitGroup
			// for e := 0; e < 10; e++ {
			// wg.Add(1)
			// go func(d int) {
			// 	c.Knock(remote, k.id)
			// 	wg.Done()
			// }(e)
			go func() {
				res, err := c.Knock(remote, k.id)
				if err == nil {
					go k.Resolve(res.Address)
				}
			}()
			// }
			select {
			case <-f:
				return
			case <-time.After(time.Second):
				// wg.Wait()
			}
		}
		x <- errors.New("knock timeout")
	}()
	select {
	case res := <-k.done:
		go func() { f <- true }()
		return res, nil
	case err := <-x:
		return "", err
	}
	// return "", nil
}

func (k Knock) RunDo(remote string, c *Client) {
	w := make(chan error)
	for i := 0; i < 5; i++ {
		// var wg sync.WaitGroup
		// for e := 0; e < 10; e++ {
		// wg.Add(1)
		// go func(d int) {
		// 	_, err := c.Knock(inc(remote, d), k.id)
		// 	wg.Done()
		// 	w <- err
		// }(e)
		// }
		go func() {
			_, err := c.Knock(remote, k.id)
			w <- err
		}()
		// wg.Wait()
		select {
		case res := <-k.done:
			log.Println("knock ok ", res)
			return
		case err := <-w:
			if err == nil {
				return
			}
		case <-time.After(time.Second):
		}
	}
}

func inc(r string, d int) string {
	rr := strings.Split(r, ":")
	x, _ := strconv.Atoi(rr[1])
	return fmt.Sprintf("%v:%v", rr[0], x+d)
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

func (p *PendingKnocksTS) Add(id string) Knock {
	p.l.Lock()
	defer p.l.Unlock()
	return p.store.Add(id)
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
