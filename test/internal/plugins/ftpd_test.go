package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestFTPClient(t *testing.T) {
	// Connect to the FTP server
	conn, err := net.Dial("tcp", "localhost:2121")
	if err != nil {
		t.Fatalf("Failed to connect to FTP server: %s", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	t.Run("Login", func(t *testing.T) {
		fmt.Println("Testing Login")

		// Read the initial server response
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read server response: %s", err)
		}

		// Check if the server is ready for login
		if !strings.HasPrefix(response, "220 ") {
			t.Fatalf("Invalid server response: %s", response)
		}

		username := "testuser"
		password := "testpass"

		// Send the username command
		cmd := fmt.Sprintf("USER %s\r\n", username)
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send username command: %s", err)
		}

		// Read the username command response
		response, err = reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read username command response: %s", err)
		}

		// Check if the username is accepted
		if !strings.HasPrefix(response, "331 ") {
			t.Fatalf("Invalid username command response: %s", response)
		}

		// Send the password command
		cmd = fmt.Sprintf("PASS %s\r\n", password)
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send password command: %s", err)
		}

		// Read the password command response
		response, err = reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read password command response: %s", err)
		}

		// Check if the password is accepted
		if !strings.HasPrefix(response, "230 ") {
			t.Fatalf("Invalid password command response: %s", response)
		}
	})

	t.Run("SYSTCommand", func(t *testing.T) {
		fmt.Println("Testing SYST Command")

		// Send the SYST command
		cmd := "SYST\r\n"
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send SYST command: %s", err)
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read SYST command response: %s", err)
		}

		// Check if the response equals "215 UNIX Type: L8."
		if response != "215 UNIX Type: L8.\r\n" {
			t.Fatalf("Invalid SYST command response: %s", response)
		}
	})
}
