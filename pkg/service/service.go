// This package provides multiple interfaces to load the services, validate them before running them
// and watching over their status
package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/riotpot/pkg/utils"
	"github.com/riotpot/pkg/validators"
)

// Interface used by every service plugin that offers a service. At the very least, every plugin
// must contain the set of methods and attributes from this interface.
// It is up to the plugin to determine the implementation of these methods for the most part.
type Service interface {
	// Get attributes from the structure
	GetID() string
	GetName() string
	GetNetwork() utils.Network
	GetInteraction() utils.Interaction
	GetPort() int
	GetAddress() string
	GetHost() string
	IsLocked() bool

	// Setters
	SetPort(port int) (int, error)
	SetName(name string)
	SetHost(host string)
	SetLocked(locked bool) (bool, error)
}

// Implements a mixin service that can be used as a base for any other service `struct` type.
type service struct {
	// require the methods described by `Service` on loading
	Service

	id          uuid.UUID
	name        string
	network     utils.Network
	port        int
	host        string
	locked      bool
	interaction utils.Interaction
}

// Getters
func (as *service) GetID() string {
	return as.id.String()
}

func (as *service) GetName() string {
	return as.name
}

func (as *service) GetNetwork() utils.Network {
	return as.network
}

func (as *service) GetInteraction() utils.Interaction {
	return as.interaction
}

func (as *service) GetPort() int {
	return as.port
}

func (as *service) GetHost() string {
	return as.host
}

func (as *service) GetAddress() string {
	return fmt.Sprintf("%s:%d", as.host, as.port)
}

func (as *service) IsLocked() bool {
	return as.locked
}

// Setters
func (as *service) SetPort(port int) (p int, err error) {
	err = validators.ValidatePortNumber(port)
	if err != nil {
		return
	}

	p = port
	as.port = port
	return
}

func (as *service) SetName(name string) {
	as.name = name
}

func (as *service) SetHost(host string) {
	as.host = host
}

func (as *service) SetLocked(locked bool) (bool, error) {
	as.locked = locked
	return as.locked, nil
}

// Implementation of a plugin-based service
// These services are stored localy as binary files that are mounted into the
// application as symbols that can be called
type PluginService interface {
	Run() error
}

type pluginService struct {
	PluginService
	Service
}

func (aps *pluginService) IsLocked() bool {
	return true
}

func (aps *pluginService) SetLocked(locked bool) (bool, error) {
	return true, fmt.Errorf("the lock status of this service can not change")
}

func NewService(name string, port int, network utils.Network, host string, interaction utils.Interaction) Service {
	return &service{
		id:          uuid.New(),
		name:        name,
		port:        port,
		network:     network,
		host:        host,
		interaction: interaction,
	}
}

// Simple constructor for plugin services
func NewPluginService(name string, port int, network utils.Network) Service {
	return &pluginService{
		Service: NewService(name, port, network, "localhost", utils.Low),
	}
}
