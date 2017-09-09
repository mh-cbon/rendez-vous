package node

import (
	"time"

	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/server"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
)

type CentralPointNode struct {
	*rendezVousNode
	cleaner       *server.Cleaner
	registrations *store.TSRegistrations
	secSocket     socket.Socket
}

func NewCentralPointNode(listen string) *CentralPointNode {
	node := &CentralPointNode{
		rendezVousNode: newrendezVousNode(listen),
		registrations:  store.NewRegistrations(nil),
	}
	node.cleaner = server.NewCleaner(time.Second*30*2, node.registrations)
	node.rendezVousNode.SetupHandler = node.setup
	return node
}
func (r *CentralPointNode) Close() error {
	r.cleaner.Stop()
	r.secSocket.Close()
	return r.rendezVousNode.Close()
}

func (r *CentralPointNode) setup(sk socket.Socket) (*server.JSONServer, *client.Client, error) {

	c := client.New(client.JSON(sk))

	socket, err := socket.FromAddr(":0")
	if err != nil {
		return nil, nil, err
	}
	err = timeout(func() error {
		return socket.Listen(nil)
	}, time.Millisecond*10)
	if err != nil {
		return nil, nil, err
	}
	r.secSocket = socket
	c2 := client.New(client.JSON(r.secSocket))

	r.cleaner.Start()

	behavior := server.OneOf(
		server.HandlePing(),
		server.HandleRegister(r.registrations),
		server.HandleUnregister(r.registrations),
		server.HandleFind(r.registrations),
		server.HandleList(r.registrations),
		server.HandleRequestKnock(r.registrations, c),
		server.HandleTestPort(r.registrations, c2),
	)
	server := server.JSON(sk, behavior)

	return server, c, nil
}
