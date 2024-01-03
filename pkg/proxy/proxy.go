package proxy

import (
	"github.com/google/uuid"
	"github.com/riotpot/pkg/service"
	"github.com/riotpot/pkg/utils"
)

// Proxy interface.
type Proxy interface {
	// Start and stop
	Start() error
	Stop()

	// Getters
	GetID() string
	GetPort() int
	GetNetwork() utils.Network
	IsRunning() utils.Status
	GetService() service.Service

	// Setters
	SetPort(port int) int
	SetService(service service.Service) service.Service
}

// Abstraction of the proxy endpoint
// Contains private fields, do not use outside of this package
type baseProxy struct {
	Proxy

	// ID of the proxy
	id uuid.UUID

	// Port in where the proxy will listen
	port int
	// Protocol meant for this proxy
	network utils.Network

	// Create a channel to stop the proxy gracefully
	// This channel is also used to guess if the proxy is running
	quit chan struct{}

	// Pointer to the slice of middlewares for the proxies
	// All the proxies should apply and share the same middlewares
	// Perhaps this can be changed in the future given the need to apply middlewares per proxy
	middlewares *middlewareManager

	// Service to proxy
	service service.Service

	// Generic listener
	listener interface{ Close() error }
}

// Function to stop the proxy from runing
func (pe *baseProxy) Stop() {
	close(pe.quit)
}

func (b *baseProxy) IsRunning() (alive utils.Status) {
	if b.quit == nil {
		return
	}

	select {
	case <-b.quit:
		return
	default:
		return utils.RunningStatus
	}
}

func (pe *baseProxy) GetID() string {
	return pe.id.String()
}

// Returns the proxy port
func (pe *baseProxy) GetPort() int {
	return pe.port
}

// Returns the service
func (pe *baseProxy) GetService() service.Service {
	return pe.service
}

// Returns the service
func (pe *baseProxy) GetNetwork() utils.Network {
	return pe.network
}

// Set the port
// NOTE: use the ValidatePort before assigning
func (pe *baseProxy) SetPort(port int) int {
	pe.port = port
	return pe.port
}

// Set the service based on the list of registered services
func (pe *baseProxy) SetService(service service.Service) service.Service {
	pe.service = service
	return pe.service
}

func newProxy(port int, network utils.Network) (px *baseProxy) {
	return &baseProxy{
		id:          uuid.New(),
		port:        port,
		network:     network,
		middlewares: Middlewares,
	}
}
