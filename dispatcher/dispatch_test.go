package dispatcher_test

import (
	"log"
	"net"
	"testing"

	"github.com/mh-cbon/rendez-vous/dispatcher"
)

func Test1(t *testing.T) {
	alice, err := net.ListenUDP("udp", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer alice.Close()
	rdispatch := dispatcher.New(alice)
	rstream1 := rdispatch.New("stream1")
	rstream2 := rdispatch.New("stream2")

	w := make(chan bool)
	go func() {
		b := make([]byte, 1024)
		n, addr, err2 := rstream1.ReadFrom(b)
		log.Printf("stream1 read: %v %v %v %q\n", n, addr, err2, string(b[:n]))
		w <- true
	}()
	go func() {
		b := make([]byte, 1024)
		n, addr, err2 := rstream2.ReadFrom(b)
		log.Printf("stream2 read: %v %v %v %q\n", n, addr, err2, string(b[:n]))
		w <- true
	}()

	bob, err := net.ListenUDP("udp", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bob.Close()
	wdispatch := dispatcher.New(bob)
	wstream0 := wdispatch.New("nop")
	wstream1 := wdispatch.New("stream1")
	wstream2 := wdispatch.New("stream2")
	s0 := []byte("The message to nop")
	n0, err0 := wstream0.WriteTo(s0, alice.LocalAddr())
	log.Println("wrote ", n0, err0)

	s1 := []byte("The message to stream1")
	n, err := wstream1.WriteTo(s1, alice.LocalAddr())
	log.Println("wrote ", n, err)

	s2 := []byte("The message to stream2")
	n2, err2 := wstream2.WriteTo(s2, alice.LocalAddr())
	log.Println("wrote ", n2, err2)

	<-w
	<-w
	//-
}
