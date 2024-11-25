package main

import (
	"fmt"
	"net"
	"time"
)

const (
	redisHost = "localhost:6379"
)

// RespCommand formats a Redis command in RESP protocol
func formatRESP(args ...string) string {
	cmd := fmt.Sprintf("*%d\r\n", len(args))
	for _, arg := range args {
		cmd += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}
	return cmd
}

// SendCommand sends a command to Redis and returns the response
func sendCommand(command string) (string, error) {
	// Connect to Redis
	conn, err := net.Dial("tcp", redisHost)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Redis: %v", err)
	}
	defer conn.Close()

	// Send command
	_, err = conn.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %v", err)
	}

	// Read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(buf[:n]), nil
}

func main() {
	fmt.Println("Testing Redis SET/GET with expiry...")

	// Test 1: Set key with expiry
	setCmd := formatRESP("SET", "test_key", "hello", "PX", "5")
	resp, err := sendCommand(setCmd)
	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
		return
	}
	fmt.Printf("SET response: %s", resp)

	// Test 2: Get key immediately
	getCmd := formatRESP("GET", "test_key")
	resp, err = sendCommand(getCmd)
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
		return
	}
	fmt.Printf("Immediate GET response: %s", resp)

	// Test 3: Get key after 3 secondsJ
	fmt.Println("Waiting 3 seconds...")
	time.Sleep(3 * time.Second)

	resp, err = sendCommand(getCmd)
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
		return
	}
	// fmt.Printf("GET after  seconds: %s", resp)

	// Test 4: Get key after expiry
	fmt.Println("Waiting 2 more seconds...")
	time.Sleep(2 * time.Second)

	resp, err = sendCommand(getCmd)
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
		return
	}
	fmt.Printf("GET after expiry: %s", resp)
}
