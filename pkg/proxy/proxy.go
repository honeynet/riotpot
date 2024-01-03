package proxy

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/riotpot/pkg/service"
	"github.com/riotpot/pkg/utils"
)

// Proxy interface.
type Proxy interface {
	// Start and stop
	Start() error
	Stop() error

	// Getters
	GetID() string
	GetPort() int
	GetNetwork() utils.Network
	GetStatus() utils.Status
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
	stop chan struct{}

	// Pointer to the slice of middlewares for the proxies
	// All the proxies should apply and share the same middlewares
	// Perhaps this can be changed in the future given the need to apply middlewares per proxy
	middlewares *middlewareManager

	// Service to proxy
	service service.Service

	// Waiting group for the server
	wg sync.WaitGroup

	// Generic listener
	listener interface{ Close() error }
}

// Function to stop the proxy from runing
func (pe *baseProxy) Stop() (err error) {
	// Stop the proxy if it is still alive
	if pe.GetStatus() == utils.RunningStatus {
		close(pe.stop)

		if pe.listener != nil {
			pe.listener.Close()
		}

		// Wait for all the connections and the server to stop
		pe.wg.Wait()
		return
	}

	err = fmt.Errorf("proxy not running")
	return
}

// Simple function to check if the proxy is running
func (pe *baseProxy) GetStatus() (alive utils.Status) {
	// When the proxy is instantiated, the stop channel is nil;
	// therefore, the proxy is not running
	if pe.stop == nil {
		return
	}

	// [7/4/2022] NOTE: The logic of this block is difficult to read.
	// However, the select block will only give the default value when there is nothing
	// to read from the channel while the channel is still open.
	// When the channel is closed, the first case is not blocked, so we can not
	// read "anything else" from the channel
	select {
	// Return if the channel is closed
	case <-pe.stop:
	// Return if the channel is open
	default:
		alive = utils.RunningStatus
	}

	return
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
		wg:          sync.WaitGroup{},
	}
}
