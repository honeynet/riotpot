// This package implements a proxy manager
// A proxy manager

package proxy

import (
	"fmt"

	"github.com/riotpot/pkg/service"
	"github.com/riotpot/pkg/utils"
)

var (
	// Instantiate the proxy manager to allow other applications work with the proxies
	Proxies = NewProxyManager()
)

// Interface for the proxy manager
type ProxyManager interface {
	// Get all the proxies registered
	GetProxies() []Proxy
	// Create a new proxy and add it to the manager
	CreateProxy(protocol string, port int) (Proxy, error)

	// Methods for proxies the using ID field
	GetProxy(id string) (Proxy, error)
	SetProxy(pe Proxy) (Proxy, error)
	DeleteProxy(id string) error

	// Wrapper method to find a proxy using the port and protocol
	GetProxyFromParams(network string, port int)

	// Set the service for a proxy
	SetService(port int, service service.Service) (pe Proxy, err error)
}

// Simple implementation of the proxy manager
// This manager has access to the proxy endpoints registered. However, it does not observe newly
type proxyManager struct {
	ProxyManager

	// List of proxy endpoints registered in the manager
	proxies []Proxy

	// Instance of the middleware manager
	middlewares *middlewareManager
}

// Create a new proxy and add it to the manager
func (pm *proxyManager) CreateProxy(network utils.Network, port int) (pe Proxy, err error) {

	// Create the proxy
	pf := &ProxyFactory{}
	pe, err = pf.CreateProxy(port, network)
	if err != nil {
		return
	}

	// Append the proxy to the list
	pm.proxies = append(pm.proxies, pe)
	return
}

func (pm *proxyManager) GetProxy(id string) (pe Proxy, err error) {
	// Get all the proxies registered
	proxies := pm.GetProxies()

	for _, proxy := range proxies {
		if proxy.GetID() == id {
			pe = proxy
			return
		}
	}

	err = fmt.Errorf("proxy not found")
	return
}

func (pm *proxyManager) SetProxy(px Proxy) (pe Proxy, err error) {
	// Get all the proxies registered
	proxies := pm.GetProxies()

	for ind, proxy := range proxies {
		if proxy.GetID() == px.GetID() {
			// Replace the index of the proxy with the new one
			proxies[ind] = px
			pm.proxies = proxies
			return
		}
	}

	err = fmt.Errorf("proxy not found")
	return
}

// Delete a proxy from teh registered list using the ID
func (pm *proxyManager) DeleteProxy(id string) (err error) {
	// Get all the proxies registered
	proxies := pm.GetProxies()

	for ind, proxy := range proxies {
		if proxy.GetID() == id {
			// Attempt to remove the service and the proxy if the service is not locked
			srv := proxy.GetService()
			if srv != nil {
				// Delete the service or return the error.
				// An error may occur if the service could not be found or is locked!
				err = service.Services.DeleteService(srv.GetID())
				if err != nil {
					return
				}
			}

			// Stop the proxy, just in case
			proxy.Stop()
			// Remove it from the slice by replacing it with the last item from the slice,
			// and reducing the slice by 1 element
			lastInd := len(proxies) - 1

			proxies[ind] = proxies[lastInd]
			pm.proxies = proxies[:lastInd]
			return
		}
	}
	return
}

func (pm *proxyManager) GetProxies() []Proxy {
	return pm.proxies
}

// Constructor for the proxy manager
func NewProxyManager() *proxyManager {
	return &proxyManager{
		middlewares: Middlewares,
		proxies:     make([]Proxy, 0),
	}
}
