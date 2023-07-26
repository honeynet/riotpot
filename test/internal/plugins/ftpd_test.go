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
	// Change port to the port running ftp server
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

	t.Run("MakeDirectory", func(t *testing.T) {
		fmt.Println("Testing Make Directory Command")

		// Send the MKD command to create the directory "test1"
		cmd := "MKD test1\r\n"
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send MKD command: %s", err)
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read MKD command response: %s", err)
		}

		// Check if the response starts with "257 " which indicates success
		if !strings.HasPrefix(response, "257 ") {
			t.Fatalf("Failed to create directory: %s", response)
		}
	})

	t.Run("ChangeDirectory", func(t *testing.T) {
		fmt.Println("Testing Change Directory Command")

		// Send the CWD command to change to the "test1" directory
		cmd := "CWD test1\r\n"
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send CWD command: %s", err)
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read CWD command response: %s", err)
		}

		// Check if the response starts with "250 " which indicates success
		if !strings.HasPrefix(response, "250 ") {
			t.Fatalf("Failed to change directory: %s", response)
		}

		// Send the CWD command to change to the "/home" directory
		cmd = "CWD /home\r\n"
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to send CWD command: %s", err)
		}

		response, err = reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read CWD command response: %s", err)
		}

		// Check if the response starts with "250 " which indicates success
		if !strings.HasPrefix(response, "250 ") {
			t.Fatalf("Failed to change directory: %s", response)
		}
	})
}
