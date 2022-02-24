package sync

import (
	"golang.org/x/xerrors"
	"time"
)

// Pool holds Clients
type Pool struct {
	connections chan interface{}
	capacity    uint
	factory     Factory
}

// Factory is a function to create new connections
type Factory func() interface{}

// NewPool creates a new pool of interface.
func NewPool(maxCap uint, factory Factory) *Pool {
	return &Pool{
		connections: make(chan interface{}, maxCap),
		capacity:    maxCap,
		factory:     factory,
	}
}

// SetFactory declare auto create function
func (p *Pool) SetFactory(factory Factory) {
	p.factory = factory
}

// GetCapacity return defined pool capacity
func (p *Pool) GetCapacity() uint {
	return p.capacity
}

func (p *Pool) Get() (rv interface{}) {
	// Try to grab an available connection within 1ms
	select {
	case rv := <-p.connections:
		return rv
	case <-time.After(time.Millisecond):
		// Try to fetch one more time, or create a new instance if factory is set else return nil
		select {
		case rv := <-p.connections:
			return rv
		default:
			if p.factory != nil {
				return p.factory()
			} else {
				return nil
			}
		}
	}
}

func (p *Pool) Put(c interface{}) error {
	select {
	case p.connections <- c:
		return nil
	default:
		return xerrors.New("pool is full, element will not be added")
	}
}
