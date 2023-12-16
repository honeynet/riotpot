package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	lr "github.com/riotpot/pkg/logger"
	"github.com/riotpot/pkg/service"
	"github.com/riotpot/pkg/utils"
)

var Plugin string

const (
	name    = "HTTPS"
	network = utils.TCP
	port    = 443
)

func init() {
	Plugin = "Httpsd"
}

type Https struct {
	// Anonymous fields from the mixin
	service.Service
}

func Httpsd() service.Service {
	mx := service.NewPluginService(name, port, network)

	return &Https{
		mx,
	}
}

func (h *Https) Run() (err error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", http.HandlerFunc(h.valid))

	cert, err := generateSelfSignedCertificate()
	if err != nil {
		lr.Log.Fatal().Err(err)
	}

	srv := &http.Server{
		Addr:    h.GetAddress(),
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}

	go h.serve(srv)

	return
}

func (h *Https) serve(srv *http.Server) {
	if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		lr.Log.Fatal().Err(err)
	}
}

// This function handles connections made to a valid path
func (h *Https) valid(w http.ResponseWriter, req *http.Request) {
	var (
		head, body string
	)

	head = `
	<html lang="en">
	<head>
		<!-- Page title -->
		<title>SCADA Login</title>

		<!-- Meta tags -->
		<meta charset="UTF-8">
		<meta id ="viewport" name="viewport" content="width=device-width, initial-scale=1">

		<!-- CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css" 
		rel="stylesheet" integrity="sha384-eOJMYsd53ii+scO/bJGFsiCZc+5NDVN2yr8+0RDqr0Ql0h+rP48ckxlpbzKgwra6" 
		crossorigin="anonymous">
	</head>
	<body>
		<h1>Login</h1><br>
	`

	body = `
		<div class="container">
			<form method="POST">
				<div class="mb-3 row">
					<label for="username" class="form-label">Username</label>
					<input id="username" name="username" class="form-control" type="text" placeholder="Username">
				</div>
				<div class="mb-3 row">
					<label for="password" class="form-label">Password</label>
					<input id="password" name="password" class="form-control" type="password" placeholder="Password">
				</div>
				<button type="submit">Log In</button>
			</form>
		</div>
	</body>
	</html>
	`

	if req.Method == http.MethodPost {
		errormessage := `
		<div class="alert alert-danger">
			<p>Incorrect username or password.</p>
		</div>
		`
		body = errormessage + body
	}

	response := fmt.Sprintf("%s%s", head, body)

	fmt.Fprint(w, response)
}

func generateSelfSignedCertificate() (tls.Certificate, error) {
	// Create certificates on the fly
	// Modified from https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// configure the certificate
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "riotpot.com",
			Organization: []string{"RiotPot"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := []string{"riotpot.com"}
	for _, h := range hosts {
		ip := net.ParseIP(h)
		if ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return tls.X509KeyPair(certPEM, keyPEM)
}
