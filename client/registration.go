package client

import (
	"log"
	"time"
)

//Registration happens at regular time intervals
type Registration struct {
	i      time.Duration
	done   chan bool
	client Client
}

// NewRegistration creates a registration for given time interval using given client
func NewRegistration(interval time.Duration, client Client) Registration {
	return Registration{interval, make(chan bool), client}
}

// Start the registration
func (r Registration) Start(remote, pbk, sign, value string) {
	go r.loop(remote, pbk, sign, value)
}

// Stop the registration
func (r Registration) Stop() {
	r.done <- true
}

func (r Registration) loop(remote, pbk, sign, value string) {
	for ; ; <-time.After(r.i) {
		_, err := r.client.Register(remote, pbk, sign, value)
		if err != nil {
			log.Println(err)
		}
		// log.Println(res)
	}
}
