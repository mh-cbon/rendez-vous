package client

import (
	"log"
	"time"

	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
)

//Registration happens at regular time intervals
type Registration struct {
	i      time.Duration
	done   chan bool
	client *Client
	remote string
	id     identity.Identity
}

// NewRegistration creates a registration for given time interval using given client
func NewRegistration(interval time.Duration, client *Client) *Registration {
	return &Registration{i: interval, done: make(chan bool), client: client}
}

// Config the registration
func (r *Registration) Config(remote string, id identity.Identity) {
	r.remote = remote
	r.id = id
}

// Start the registration
func (r *Registration) Start() error {
	go r.loop()
	return nil
}

// Stop the registration
func (r *Registration) Stop() error {
	r.done <- true
	return nil
}

func (r *Registration) register() {
	res, err := r.client.Register(r.remote, &r.id)
	if err != nil {
		log.Println(err)
	}
	if res.Code != model.OkCode {
		log.Println("registration failed", res)
	}
}

func (r *Registration) loop() {
	go r.register()
	for {
		select {
		case <-r.done:
			return
		case <-time.After(r.i):
			go r.register()
		}
	}
}
