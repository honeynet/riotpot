package main

import (
	"io"
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

func TestUpnpdServer(t *testing.T) {
	// Create an HTTP client
	client := &http.Client{}

	// Create an M-POST request
	req, err := http.NewRequest("M-POST", "http://localhost:1900", nil)
	assert.NoError(t, err)
	req.Header.Set("SOAPAction", "GetExternalIPAddress")

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check the response content type
	assert.Equal(t, "text/xml", resp.Header.Get("Content-Type"))

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Check the response body
	expectedResponse := `<?xml version="1.0"?>
			<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Body>
					<u:GetExternalIPAddressResponse xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
						<NewExternalIPAddress>192.168.1.100</NewExternalIPAddress>
					</u:GetExternalIPAddressResponse>
				</s:Body>
			</s:Envelope>
		`

	expectedResponse = removeWhitespace(expectedResponse)
	bodyString := removeWhitespace(string(body))
	assert.Equal(t, expectedResponse, bodyString)
}
