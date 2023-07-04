package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeWhitespace(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func TestValidPathHandler(t *testing.T) {
	// Disable certificate verification for testing purposes
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Make an HTTPS request to the default HTTPS port (443)
	resp, err := client.Get("https://localhost:443")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Assert that the response body contains the expected HTML content
	expectedBody := `<html lang="en">
	<head>
		<!-- Page title -->
		<title>SCADA Login</title>

		<!-- Meta tags -->
		<meta charset="UTF-8">
		<meta id="viewport" name="viewport" content="width=device-width, initial-scale=1">

		<!-- CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css" 
		rel="stylesheet" integrity="sha384-eOJMYsd53ii+scO/bJGFsiCZc+5NDVN2yr8+0RDqr0Ql0h+rP48ckxlpbzKgwra6" 
		crossorigin="anonymous">
	</head>
	<body>
		<h1>Login</h1><br>
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
	</html>`

	// Remove newline and tab characters from expectedBody and body
	expectedBody = removeWhitespace(expectedBody)
	bodyStr := removeWhitespace(string(body))
	assert.Equal(t, expectedBody, bodyStr)

	// Assert the TLS certificate
	connState := resp.TLS
	assert.NotNil(t, connState)

	// Get the first certificate in the chain
	certs := connState.PeerCertificates
	assert.NotEmpty(t, certs)
	cert := certs[0]

	// Assert the subject common name of the certificate
	expectedCN := "riotpot.com"
	assert.Equal(t, expectedCN, cert.Subject.CommonName)
}
