package proxy

import (
	"fmt"

	"github.com/riotpot/pkg/utils"
)

type ProxyFactory struct{}

func (pfac *ProxyFactory) CreateProxy(port int, network utils.Network) (px Proxy, err error) {
	switch network {
	case utils.TCP:
		px, err = NewTCPProxy(port)
	case utils.UDP:
		px, err = NewUDPProxy(port)
	default:
		return nil, fmt.Errorf("proxy not found")
	}

	return
}
