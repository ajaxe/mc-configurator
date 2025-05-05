package main

import (
	"os"
	"testing"
)

func TestEnvConfigRead(t *testing.T) {
	// Set the environment variable for testing
	os.Setenv("MC_LEVEL_NAME", "test_level")

	// Call the LoadConfig function
	config := loadConfig()

	// Check if the config contains the expected value
	if config["level-name"] != "test_level" {
		t.Errorf("Expected level-name to be 'test_level', got '%s'", config["level-name"])
	}

	v := os.Getenv("MC_LEVEL_NAME")

	if v != "test_level" {
		t.Errorf("Expected level-name to be 'test_level', got '%s'", config["level-name"])
	}

	// Clean up the environment variable
	os.Unsetenv("MC_LEVEL_NAME")
}
