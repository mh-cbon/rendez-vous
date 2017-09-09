package server

import (
	"log"
	"time"

	"github.com/mh-cbon/rendez-vous/store"
)

func NewCleaner(ttl time.Duration, store *store.TSRegistrations) *Cleaner {
	return &Cleaner{ttl: ttl, store: store}
}

type Cleaner struct {
	store *store.TSRegistrations
	done  chan bool
	ttl   time.Duration
}

func (c *Cleaner) Start() error {
	c.done = make(chan bool)
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-time.After(time.Second * 2):
				c.store.Transact(func(store *store.Registrations) {
					for i, p := range store.Peers {
						if p.Create.Add(c.ttl).Before(time.Now()) {
							log.Println("Clean ", p)
							log.Println("Clean ", p)
							if i+1 < len(store.Peers) {
								store.Peers = append(store.Peers[:i], store.Peers[i+1:]...)
							} else {
								store.Peers = append(store.Peers[:0], store.Peers[:i]...)
							}
						}
					}
				})
			}
		}
	}()
	return nil
}
func (c *Cleaner) Stop() error {
	c.done <- true
	return nil
}
