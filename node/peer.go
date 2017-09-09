package node

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/identity"
	"github.com/mh-cbon/rendez-vous/model"
	"github.com/mh-cbon/rendez-vous/server"
	"github.com/mh-cbon/rendez-vous/socket"
	"github.com/mh-cbon/rendez-vous/store"
)

type PeerNode struct {
	*rendezVousNode
	Port      *PortStatus
	knocks    *store.TSPendingOps
	portTests *store.TSPendingOps
}

func NewPeerNode(listen string) *PeerNode {
	node := &PeerNode{
		rendezVousNode: newrendezVousNode(listen),
		Port:           NewPortStatus(),
		knocks:         store.NewTSPendingOps(nil),
		portTests:      store.NewTSPendingOps(nil),
	}
	node.rendezVousNode.SetupHandler = node.setup
	return node
}

func (r *PeerNode) setup(sk socket.Socket) (*server.JSONServer, *client.Client, error) {

	client := client.New(client.JSON(sk))

	behavior := server.OneOf(
		server.HandlePing(),
		server.HandleDoKnock(client, r.knocks),
		server.HandleKnock(r.knocks),
		server.HandlePortTest(r.portTests),
	)
	server := server.JSON(sk, behavior)

	r.Port.Set(model.PortStatusUnknown, sk.LocalAddr().(*net.UDPAddr).Port)

	return server, client, nil
}

func (r *PeerNode) ReqKnock(remote string, id *identity.PublicIdentity) (net.Addr, error) {
	token := "random"
	knock := r.knocks.Add(token, func(remote string, m model.Message) bool {
		return m.Token == token
	})
	defer r.knocks.Rm(knock)
	c := r.GetClient()
	f, err := c.ReqKnock(remote, id, token)
	if err == nil {
		for i := 0; i < 5; i++ {
			res, err2 := knock.Run(func() error {
				go c.Knock(f.Data, token)
				return nil
			})
			if err2 == nil {
				addr, err3 := net.ResolveUDPAddr("udp", res.Remote)
				if err3 == nil {
					return addr, nil
				}
			}
		}
		return nil, errors.New("knock failed")
	}
	return nil, err
}

func (r *PeerNode) TestPort(remote string, h func(int, int)) *PortStatus {
	token := "random"
	client := r.GetClient()
	portTest := r.portTests.Add(token, func(remote string, m model.Message) bool {
		return m.Token == token
	})
	defer r.portTests.Rm(portTest)
	res, err := portTest.Run(func() error {
		_, err := client.TestPort(remote, token)
		return err
	})
	if err == nil {
		a, err := net.ResolveUDPAddr("udp", res.M.Address)
		if err == nil {
			r.Port.Set(model.PortStatusOpen, a.Port)
		}
	} else {
		r.Port.Set(model.PortStatusClose, r.LocalAddr().(*net.UDPAddr).Port)
	}
	return r.Port
}

func (r *PeerNode) Resolve(remote string, addr string, service string, me *identity.PublicIdentity) (string, error) {
	h := strings.Split(addr, ":")
	host := h[0]
	if strings.HasSuffix(host, ".me.com") {
		pbk := host[:len(host)-7]
		remoteID, err2 := identity.FromPbk(pbk, service)
		if err2 != nil {
			return "", err2
		}
		c := r.GetClient()
		res, err2 := c.Find(remote, remoteID)
		if err2 != nil {
			return "", err2
		}
		if res.PortStatus == model.PortStatusClose {
			if me == nil {
				return "", errors.New("failed to request knock: me identity is missing")
			}
			newRemote, err2 := r.ReqKnock(remote, me)
			log.Println("found ", newRemote.String())
			log.Println("err2 ", err2)
			if err2 != nil {
				return "", fmt.Errorf("knock failure: %v", err2.Error())
			}
			return newRemote.String(), nil
		}
		log.Printf("%#v\n", res)
		return res.Data, nil
	}
	return addr, nil
}

func (r *PeerNode) Dial(addr string, svc string) (net.Conn, error) {
	id := []byte(fmt.Sprintf("%v:%v", len(svc), svc))
	log.Println("addr", addr)
	conn, err := r.PacketConn().Dial(addr)
	log.Println("conn", err)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(id)
	log.Println("Write", err)
	if err != nil {
		conn.Close()
		return nil, err
	}
	b := make([]byte, len(id))
	n, err := conn.Read(b)
	log.Println("Read", err)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if bytes.Equal(id, b[:n]) == false {
		conn.Close()
		return nil, fmt.Errorf("Invalid service %v", string(b))
	}
	return conn, nil
}

func (r *PeerNode) ServiceListener(svc string) net.Listener {
	return &serviceListener{
		Listener: r.ln,
		id:       []byte(fmt.Sprintf("%v:%v", len(svc), svc)),
	}
}

type serviceListener struct {
	net.Listener
	id []byte
}

func (s *serviceListener) Accept() (net.Conn, error) {
	conn, err := s.Listener.Accept()
	if err == nil {
		b := make([]byte, len(s.id))
		n, err2 := conn.Read(b)
		if err != nil {
			conn.Close()
			return nil, err2
		}
		if bytes.Equal(s.id, b[:n]) == false {
			conn.Close()
			return nil, fmt.Errorf("Invalid service %v", string(b))
		}
		_, err2 = conn.Write(s.id)
		if err2 != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, err
}
