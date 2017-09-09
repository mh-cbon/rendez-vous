package node

import (
	"sync"

	"github.com/mh-cbon/rendez-vous/model"
)

type PortStatus struct {
	l        sync.Mutex
	num      int
	status   int
	onchange func(newStatus int, oldStatus int)
}

func NewPortStatus() *PortStatus {
	return &PortStatus{
		l:      sync.Mutex{},
		num:    0,
		status: model.PortStatusUnknown,
	}
}

func (p *PortStatus) OnChange(h func(int, int)) *PortStatus {
	p.l.Lock()
	defer p.l.Unlock()
	p.onchange = h
	return p
}

func (p *PortStatus) Set(status, num int) bool {
	if status != p.status {
		p.l.Lock()
		defer p.l.Unlock()
		old := p.status
		p.num = num
		p.status = status
		if p.onchange != nil {
			p.onchange(status, old)
		}
		return true
	}
	return false
}

func (p *PortStatus) Num() int {
	p.l.Lock()
	defer p.l.Unlock()
	return p.num
}
func (p *PortStatus) Status() int {
	p.l.Lock()
	defer p.l.Unlock()
	return p.status
}
func (p *PortStatus) Open() bool {
	return p.Status() == model.PortStatusOpen
}
func (p *PortStatus) Close() bool {
	return p.Status() == model.PortStatusClose
}
