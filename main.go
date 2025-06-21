// This program updates the server.properties file for a Minecraft server with new values from environment variables.
// Or use 'rcon' to send commands to Minecraft server.
package main

import (
	"fmt"
	"mc-configurator/internal/config"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage(os.Args[0])
	}
	if strings.ToLower(os.Args[1]) != "rcon" {
		o := config.NewGeneratorOptions(os.Args[1:])
		generator := config.NewConfigGenerator(o)
		generator.Execute()
	}
}

func printUsage(n string) {
	_, f := filepath.Split(n)
	fmt.Printf("Usage: %s [rcon <subcommand> |<src file> <dest file>]\n", f)
	println("This program updates the server.properties file for a Minecraft server with new values from environment variables. Or use 'rcon' to send commands to Minecraft server.")
	println("Options:")
	println("  rcon <subcommand> - Executes <subcommand> via RCON to the Minecraft server.")
	println("  <src file> <dest file> - Using <src> file as template applies environemt variable as overrides generates new  server properties at <dest> file path.")
	os.Exit(1)
}
