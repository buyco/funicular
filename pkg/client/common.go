package client

import (
	"github.com/buyco/keel/pkg/helper"
	"github.com/sirupsen/logrus"
	"time"
)

// Pool holds Clients
type Pool struct {
	connections chan interface{}
	factory     Factory
	logger      *logrus.Logger
}

// Factory is a function to create new connections
type Factory func() interface{}

// NewPool creates a new pool of interface.
func NewPool(maxCap uint, factory Factory, logger *logrus.Logger) *Pool {
	return &Pool{
		connections: make(chan interface{}, maxCap),
		factory:     factory,
		logger:      logger,
	}
}

// SetFactory declare auto create function
func (p *Pool) SetFactory(factory Factory) {
	p.factory = factory
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
		return helper.ErrorPrint("Pool is full, element will not be added")
	}
}
