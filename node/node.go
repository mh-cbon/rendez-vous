package node

//
// import (
// 	"errors"
// 	"net"
// 	"sync"
// 	"time"
//
// 	"github.com/anacrolix/utp"
// 	"github.com/mh-cbon/rendez-vous/client"
// 	"github.com/mh-cbon/rendez-vous/server"
// 	"github.com/mh-cbon/rendez-vous/socket"
// )
//
// func newrendezVousNode(listen string) *rendezVousNode {
// 	return &rendezVousNode{
// 		l:      &sync.Mutex{},
// 		listen: listen,
// 	}
// }
//
// type rendezVousNode struct {
// 	l            *sync.Mutex
// 	ln           net.Listener
// 	listen       string
// 	socket       socket.Socket
// 	server       *server.JSONServer
// 	client       *client.Client
// 	SetupHandler func(socket socket.Socket) (*server.JSONServer, *client.Client, error)
// }
//
// func (r *rendezVousNode) GetClient() *client.Client {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.client
// }
//
// func (r *rendezVousNode) Listener() net.Listener {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.ln
// }
//
// func (r *rendezVousNode) PacketConn() *utp.Socket {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.ln.(*utp.Socket)
// }
//
// func (r *rendezVousNode) Start() error {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	ln, err := utp.Listen(r.listen)
// 	if err != nil {
// 		return err
// 	}
// 	r.ln = ln
// 	r.socket = socket.FromConn(r.ln.(*utp.Socket))
//
// 	server, client, err := r.SetupHandler(r.socket)
// 	if err != nil {
// 		return err
// 	}
// 	r.client = client
// 	r.server = server
// 	return timeout(r.server.ListenAndServe, time.Millisecond*10)
// }
//
// func (r *rendezVousNode) Close() error {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.socket.Close()
// }
//
// func (r *rendezVousNode) Listen() string {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.listen
// }
//
// func (r *rendezVousNode) LocalAddr() net.Addr {
// 	r.l.Lock()
// 	defer r.l.Unlock()
// 	return r.socket.LocalAddr()
// }
//
// func (r *rendezVousNode) Restart(listen string) error {
// 	if r.listen == listen {
// 		return errors.New("must be different")
// 	}
// 	if err := r.Close(); err != nil {
// 		return err
// 	}
// 	r.l.Lock()
// 	r.listen = listen
// 	r.l.Unlock()
// 	return r.Start()
// }
//
// func timeout(do func() error, d time.Duration) error {
// 	rcv := make(chan error)
// 	go func() {
// 		rcv <- do()
// 	}()
// 	select {
// 	case err := <-rcv:
// 		close(rcv)
// 		if err != nil {
// 			return err
// 		}
// 	case <-time.After(d):
// 	}
// 	return nil
// }
