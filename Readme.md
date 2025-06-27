# Minecraft Server Configurator

A simple command-line tool to configure a Minecraft Bedrock `server.properties` file using environment variables. It's designed to be used within a Docker environment to dynamically set server properties at container startup.

## Overview

When running a Minecraft server in Docker, managing the `server.properties` file can be cumbersome. This tool simplifies the process by reading environment variables, translating them into the correct property format, and updating the `server.properties` file before the main server process starts.

This allows you to manage your server's configuration directly from your `docker-compose.yml` or Docker run command, making it easy to version control and replicate server setups.

## How It Works

The tool performs the following steps:

1.  Reads all environment variables from the current environment.
2.  Identifies variables that start with the `MC_` prefix.
3.  For each matching variable, it converts the variable name into a `server.properties` key (e.g., `MC_LEVEL_NAME` becomes `level-name`).
4.  It reads the existing `server.properties` file.
5.  It updates the value for each corresponding key in the properties file. If the key doesn't exist, it's appended to the file.
6.  The updated properties are written back to the `server.properties` file.

## Usage

### Environment Variables

To override a property, set an environment variable with the `MC_` prefix. The rest of the variable name should match the property key, with any dashes (`-`) replaced by underscores (`_`).

- **Prefix:** `MC_`
- **Key format:** `PROPERTY_KEY_WITH_UNDERSCORES`
- **Example:** To set `level-name`, use the environment variable `MC_LEVEL_NAME`.

### Example Overrides

Here are some common properties you can set:

```bash
# Sets 'level-name=My Awesome World' in server.properties
MC_LEVEL_NAME="My Awesome World"

# Sets 'level-seed=123456789'
MC_LEVEL_SEED="123456789"

# Sets 'gamemode=creative'
MC_GAMEMODE="creative"

# Sets 'difficulty=hard'
MC_DIFFICULTY="hard"

# Sets 'max-players=20'
MC_MAX_PLAYERS="20"

# Enables RCON and sets the password
MC_ENABLE_RCON="true"
MC_RCON_PASSWORD="your-secure-password"
```

### Docker Integration

This tool is most effective when used as an entrypoint or initial command in a Docker container. You would run `mc-configurator` first, and then start the Minecraft server.

Here is an example snippet for a `docker-compose.yml` file:

```yaml
services:
  minecraft:
    image: itzg/minecraft-bedrock-server
    ports:
      - "19132:19132/udp"
    environment:
      - EULA=TRUE
      - MC_LEVEL_NAME=My Docker World
      - MC_GAMEMODE=survival
      - MC_DIFFICULTY=normal
    volumes:
      - ./mc-data:/data
      - ./mc-configurator:/configurator # Mount the configurator binary
    command: >
      bash -c "/configurator/mc-configurator && /usr/local/bin/start"
      # 1. Run the configurator to update server.properties
      # 2. Run the original server start script
```

_Note: The exact `command` will depend on the base Minecraft Docker image you are using._

## Building

To build the binary from the source:

```bash
go build -o mc-configurator .
```

## Debug/Test

To test the configurator locally on Windows using PowerShell:

```powershell
$env:MC_LEVEL_NAME="test world"; .\mc-configurator.exe; $env:MC_LEVEL_NAME=$null
```
