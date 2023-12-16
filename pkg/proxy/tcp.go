package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	lr "github.com/riotpot/pkg/logger"
	"github.com/riotpot/pkg/utils"
	"github.com/riotpot/pkg/validators"
)

// Implementation of a TCP proxy

type tcpProxy struct {
	*baseProxy
	listener net.Listener
}

// Start listening for connections
func (px *tcpProxy) Start() (err error) {
	// Check if the service is set, otherwise return with an error
	if px.GetService() == nil {
		return fmt.Errorf("service not set")
	}

	// Get the listener or create a new one
	listener, err := px.GetListener()

	if err != nil {
		return
	}

	// Add a waiting task
	px.wg.Add(1)

	go func() {
		defer px.wg.Done()

		for {
			// Accept the next connection
			// This goes first as it is the method we have to check if the proxy is running
			// There is no need to continue if it is not
			client, err := listener.Accept()
			if err != nil {
				return
			}
			defer client.Close()

			// Apply the middlewares to the connection before dialing the server
			_, err = px.middlewares.Apply(client)
			if err != nil {
				return
			}

			// Get a connection to the server for each new connection with the client
			server, servErr := net.DialTimeout(utils.TCP.String(), px.service.GetAddress(), 1*time.Second)

			// If there was an error, close the connection to the server and return
			if servErr != nil {
				return
			}
			defer server.Close()

			// Add a waiting task
			px.wg.Add(1)

			go func() {
				// Handle the connection between the client and the server
				// NOTE: The handlers will defer the connections
				px.handle(client, server)

				// Finish the task
				px.wg.Done()
			}()
		}
	}()

	return
}

func (px *tcpProxy) GetListener() (listener net.Listener, err error) {
	listener = px.listener

	// Get the listener only
	if listener == nil || px.GetStatus() != utils.RunningStatus {
		listener, err = px.NewListener()
		if err != nil {
			return
		}
		px.listener = listener
	}

	return
}

func (px *tcpProxy) SafeSetPort(port int) (p int, err error) {
	p, err = validators.ValidatePort(port)
	if err != nil {
		return
	}

	px.SetPort(p)
	return
}

func (px *tcpProxy) NewListener() (listener net.Listener, err error) {
	listener, err = net.Listen(px.GetNetwork().String(), fmt.Sprintf(":%d", px.GetPort()))
	px.listener = listener
	return
}

// TCP synchronous tunnel that forwards requests from source to destination and back
func (px *tcpProxy) handle(from net.Conn, to net.Conn) {
	// Create the waiting group for the connections so they can answer the each other
	var wg sync.WaitGroup
	wg.Add(2)

	handler := func(source net.Conn, dest net.Conn) {
		defer wg.Done()

		// Write the content from the source to the destination
		_, err := io.Copy(dest, source)
		if err != nil {
			lr.Log.Warn().Err(err).Msg("Could not copy from source to destination")
		}

		// Close the connection to the source
		if err := source.Close(); err != nil {
			lr.Log.Warn().Err(err)
		}

		// Attempt to close the writter. This may not always work
		// Another solution is to just call `Close()` on the writter
		if d, ok := dest.(*net.TCPConn); ok {
			if err := d.CloseWrite(); err != nil {
				lr.Log.Warn().Err(err)
			}

		}
	}

	// Start the workers
	// TODO: [7/3/2022] Check somewhere if the connection is still alive from the source and destination
	// Otherwise there is no need to wait
	go handler(from, to)
	go handler(to, from)

	// Wait until the forwarding is done
	wg.Wait()
}

func NewTCPProxy(port int) (proxy *tcpProxy, err error) {
	// Create a new proxy
	proxy = &tcpProxy{
		baseProxy: newProxy(port, utils.TCP),
	}

	// Set the port
	proxy.SafeSetPort(port)
	return
}
