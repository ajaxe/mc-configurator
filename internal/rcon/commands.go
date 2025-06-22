package rcon

import (
	"fmt"
	"strings"
)

func NewCommand() *RCONCommand {
	return &RCONCommand{
		client: NewClient(),
	}
}

type RCONCommand struct {
	client *Client
}

func (r *RCONCommand) Execute(args []string) (err error) {
	cmd := strings.Join(args, " ")
	// Placeholder for command execution logic
	// This could involve sending a command to the RCON server and handling the response

	err = r.client.Connect()
	if err != nil {
		return
	}
	defer r.client.Disconnect()

	err = r.client.Authenticate()
	if err != nil {
		return
	}

	response, err := r.client.ExecuteCommand(cmd)
	if err != nil {
		fmt.Printf("Error executing command: %s: %v\n", cmd, err)
		return
	}

	fmt.Printf("Server response:\n%s\n", response)

	return nil
}
