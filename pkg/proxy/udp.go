package proxy

import (
	"fmt"
	"net"
	"sync"

	lr "github.com/riotpot/pkg/logger"
	"github.com/riotpot/pkg/utils"
)

type udpProxy struct {
	*baseProxy
	listener *net.UDPConn
}

func (px *udpProxy) Start() (err error) {
	// Check if the service is set
	if px.GetService() == nil {
		err = fmt.Errorf("service not set")
		return
	}

	// Create a channel to stop the proxy
	px.quit = make(chan struct{})

	// Get the listener or create a new one
	client, err := px.GetListener()
	if err != nil {
		return
	}
	defer client.Close()

	// Add a waiting task
	var wg sync.WaitGroup
	wg.Add(1)

	srvAddr := net.UDPAddr{
		Port: px.service.GetPort(),
	}

	go func() {
		defer wg.Done()

		for {
			select {
			case <-px.quit:
				return

			default:
				// Get a connection to the server for each new connection with the client
				server, servErr := net.DialUDP(utils.UDP.String(), nil, &srvAddr)
				// If there was an error, close the connection to the server and return
				if servErr != nil {
					return
				}
				defer server.Close()

				// Add a waiting task
				wg.Add(1)

				go func() {
					defer wg.Done()
					// TODO: Handle the middlewares! they only accept TCP connections
					// Apply the middlewares to the connection
					//udpProxy.middlewares.Apply(listener)

					// Handle the connection between the client and the server
					// NOTE: The handlers will defer the connections
					px.handle(client, server)
				}()
			}
		}
	}()

	return
}

// Get or create a new listener
func (px *udpProxy) GetListener() (listener *net.UDPConn, err error) {
	listener = px.listener

	// Check if there is a listener
	if listener == nil || px.IsRunning() != utils.RunningStatus {
		// Get the address of the UDP server
		addr := net.UDPAddr{
			Port: px.service.GetPort(),
		}

		listener, err = net.ListenUDP(utils.UDP.String(), &addr)
		if err != nil {
			return
		}
		px.listener = listener
		px.baseProxy.listener = listener
	}

	return
}

// TODO: Test this function
// UDP asynchronous tunnel
func (px *udpProxy) handle(client *net.UDPConn, server *net.UDPConn) {
	var buf [2 << 10]byte
	var wg sync.WaitGroup
	wg.Add(2)

	// Function to copy messages from one pipe to the other
	var handle = func(from *net.UDPConn, to *net.UDPConn) {
		n, addr, err := from.ReadFrom(buf[0:])
		if err != nil {
			lr.Log.Warn().Err(err)
		}

		_, err = to.WriteTo(buf[:n], addr)
		if err != nil {
			lr.Log.Warn().Err(err)
		}
	}

	defer client.Close()
	defer server.Close()

	go handle(client, server)
	go handle(server, client)

	// Wait until the forwarding is done
	wg.Wait()
}

func NewUDPProxy(port int) (proxy *udpProxy, err error) {
	// Create a new proxy
	proxy = &udpProxy{
		baseProxy: newProxy(port, utils.UDP),
	}

	// Set the port
	proxy.SetPort(port)
	return
}
