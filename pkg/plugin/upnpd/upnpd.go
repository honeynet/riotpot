package main

import (
	"fmt"
	"net/http"
	"strings"
	

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string

const (
	name    = "UPNP"
	network = globals.TCP
	port    = 1900
)

func init() {
	Plugin = "Upnpd"
}

func Upnpd() services.Service {
	mx := services.NewPluginService(name, port, network)

	return &Upnp{
		mx,
	}
}

type Upnp struct {
	// Anonymous fields from the mixin
	services.Service
}

func (h *Upnp) Run() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.valid))

	srv := &http.Server{
		Addr:    h.GetAddress(),
		Handler: mux,
	}

	go h.serve(srv)

	return
}

func (h *Upnp) serve(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lr.Log.Fatal().Err(err)
	}
}

func (h *Upnp) valid(w http.ResponseWriter, req *http.Request) {
	if req.Method == "M-POST" && strings.Contains(req.Header.Get("SOAPAction"), "GetExternalIPAddress") {
		response := `
			<?xml version="1.0"?>
			<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Body>
					<u:GetExternalIPAddressResponse xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
						<NewExternalIPAddress>192.168.1.100</NewExternalIPAddress>
					</u:GetExternalIPAddressResponse>
				</s:Body>
			</s:Envelope>
		`
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
