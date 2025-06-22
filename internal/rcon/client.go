package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	// Packet Types from the specification
	SERVERDATA_AUTH           int32 = 3
	SERVERDATA_AUTH_RESPONSE  int32 = 2
	SERVERDATA_EXECCOMMAND    int32 = 2
	SERVERDATA_RESPONSE_VALUE int32 = 0
)

// RCONPacket defines the structure of an RCON packet as per the specification.
type RCONPacket struct {
	Size      int32
	RequestID int32
	Type      int32
	Payload   []byte
	// The specification notes a 1-byte null padding, which we handle during packing/unpacking
}

// Client represents an RCON client.
type Client struct {
	conn          net.Conn
	host          string
	port          int
	password      string
	authenticated bool
	requestID     int32
}

// NewClient creates a new RCON client.
// host: The IP address or hostname of the Minecraft server.
// port: The RCON port of the server (default is 25575).
// password: The RCON password.
func NewClient() *Client {
	host := "localhost" // or your server's IP
	port := os.Getenv("MC_RCON_PORT")
	password := os.Getenv("MC_RCON_PASSWORD")
	p, _ := strconv.Atoi(port)

	fmt.Printf("Client: host=%s port=%d password=%v\n", host, p, password != "")

	return &Client{
		host:     host,
		port:     p,
		password: password,
		// Start with a non-zero request ID, for example, the current time
		requestID: int32(time.Now().Unix()),
	}
}

// Connect establishes a TCP connection to the RCON server.
func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to RCON server: %w", err)
	}
	c.conn = conn
	return nil
}

// Disconnect closes the TCP connection.
func (c *Client) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
		c.authenticated = false
	}
}

// Authenticate sends an authentication request to the server.
func (c *Client) Authenticate() error {
	if c.conn == nil {
		return errors.New("not connected to the server")
	}

	// Send authentication packet
	err := c.writePacket(SERVERDATA_AUTH, []byte(c.password))
	if err != nil {
		return fmt.Errorf("failed to send auth packet: %w", err)
	}

	// Read response
	resp, err := c.readPacket()
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	// Check for authentication failure (RequestID == -1)
	if resp.RequestID == -1 {
		c.authenticated = false
		return errors.New("authentication failed: invalid password")
	}

	// Check if the response ID matches the request ID
	if resp.RequestID != c.requestID {
		c.authenticated = false
		return fmt.Errorf("authentication failed: mismatched request IDs (expected %d, got %d)", c.requestID, resp.RequestID)
	}

	c.authenticated = true
	return nil
}

// ExecuteCommand sends a command to the server and returns the response.
// This must be called after a successful authentication.
func (c *Client) ExecuteCommand(command string) (string, error) {
	if !c.authenticated {
		return "", errors.New("not authenticated")
	}

	// Send command packet
	err := c.writePacket(SERVERDATA_EXECCOMMAND, []byte(command))
	if err != nil {
		return "", fmt.Errorf("failed to send command packet: %w", err)
	}

	// Read response
	resp, err := c.readPacket()
	if err != nil {
		return "", fmt.Errorf("failed to read command response: %w", err)
	}

	return string(resp.Payload), nil
}

// writePacket constructs and sends a packet to the server.
func (c *Client) writePacket(packetType int32, payload []byte) error {
	c.requestID++ // Increment request ID for each new packet sent

	// As per spec: Length = 4 (reqID) + 4 (type) + len(payload) + 1 (null) + 1 (null)
	// The formula `10 + len(payload)` from the spec is slightly simplified.
	// It's 4 (ID) + 4 (Type) + len(payload) + 2 (two null terminators).
	// Size = 4(ID) + 4(Type) + len(body) + 2(null bytes)
	size := int32(4 + 4 + len(payload) + 2)

	buf := new(bytes.Buffer)
	// Write Size
	binary.Write(buf, binary.LittleEndian, size)
	// Write Request ID
	binary.Write(buf, binary.LittleEndian, c.requestID)
	// Write Type
	binary.Write(buf, binary.LittleEndian, packetType)
	// Write Payload
	buf.Write(payload)
	// Write two null bytes (payload terminator and padding)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)

	_, err := c.conn.Write(buf.Bytes())
	return err
}

// readPacket reads and parses a packet from the server.
func (c *Client) readPacket() (*RCONPacket, error) {
	// Read the packet size (first 4 bytes)
	var size int32
	if err := binary.Read(c.conn, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("error reading packet size: %w", err)
	}

	// Read the rest of the packet
	data := make([]byte, size)
	if _, err := c.conn.Read(data); err != nil {
		return nil, fmt.Errorf("error reading packet data: %w", err)
	}

	// Create a reader for the packet data
	reader := bytes.NewReader(data)

	// Read RequestID
	var requestID int32
	binary.Read(reader, binary.LittleEndian, &requestID)

	// Read Type
	var packetType int32
	binary.Read(reader, binary.LittleEndian, &packetType)

	// The rest is payload, minus the two null bytes at the end
	payload := make([]byte, reader.Len()-2)
	reader.Read(payload)

	return &RCONPacket{
		Size:      size,
		RequestID: requestID,
		Type:      packetType,
		Payload:   payload,
	}, nil
}

/*
// main function to demonstrate the RCON client.
func main() {
	// --- IMPORTANT ---
	// Replace with your server's details.
	// Ensure RCON is enabled in your server.properties file.
	// enable-rcon=true
	// rcon.port=25575
	// rcon.password=your_password

	// Create a new client
	client := NewClient(serverHost, serverPort, rconPassword)

	// Connect to the server
	fmt.Println("Connecting to RCON server...")
	err := client.Connect()
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	defer client.Disconnect()
	fmt.Println("Connected.")

	// Authenticate
	fmt.Println("Authenticating...")
	err = client.Authenticate()
	if err != nil {
		fmt.Printf("Error authenticating: %v\n", err)
		return
	}
	fmt.Println("Authenticated successfully.")

	// Execute a command
	command := "list"
	fmt.Printf("\nExecuting command: '%s'\n", command)
	response, err := client.ExecuteCommand(command)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}
	fmt.Printf("Server response:\n%s\n", response)

	// Execute another command
	command = "say Hello from RCON!"
	fmt.Printf("\nExecuting command: '%s'\n", command)
	response, err = client.ExecuteCommand(command)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}
	// The "say" command doesn't return a visible response payload, but it will appear in the server chat.
	fmt.Println("Command sent. Check the in-game chat or server console.")
} */
