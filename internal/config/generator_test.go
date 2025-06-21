package config

import (
	"os"
	"testing"
)

func TestEnvConfigRead(t *testing.T) {
	levelName := "test_level"
	enableRCON := "true"
	rconPort := "12345"
	rconPassword := "test123"
	// Set the environment variable for testing
	os.Setenv("MC_LEVEL_NAME", levelName)
	os.Setenv("MC_ENABLE_RCON", enableRCON)
	os.Setenv("MC_RCON_PORT", rconPort)
	os.Setenv("MC_RCON_PASSWORD", rconPassword)

	// Call the LoadConfig function
	config := loadEnvConfig()

	// Check if the config contains the expected value
	if config["level-name"] != "test_level" {
		t.Errorf("Expected level-name to be 'test_level', got '%s'", config["level-name"])
	}

	if levelName != "test_level" {
		t.Errorf("Expected level-name to be 'test_level', got '%s'", config["level-name"])
	}
	if config["enable-rcon"] != enableRCON {
		t.Errorf("Expected enable-rcon to be 'true', got '%s'", config["enable-rcon"])
	}
	if config["rcon.port"] != rconPort {
		t.Errorf("Expected rcon.port to be '12345', got '%s'", config["rcon.port"])
	}
	if config["rcon.password"] != rconPassword {
		t.Errorf("Expected rcon.password to be 'test123', got '%s'", config["rcon.password"])
	}

	// Clean up the environment variable
	os.Unsetenv("MC_LEVEL_NAME")
	os.Unsetenv("MC_ENABLE_RCON")
	os.Unsetenv("MC_RCON_PORT")
	os.Unsetenv("MC_RCON_PASSWORD")
}
