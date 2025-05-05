# Minecraft Server Configurator

Simple configuration tool which overrides _server.properties_ file using the values set as environment variables. This allows us to change Minecraft server properties by simply updating _docker-compose_ file.

Environment vairable must have `MC_` as prefix. Property names with a `-` (dash) must be set as `_` (underscore).

Sample overrides

```bash
MC_LEVEL_NAME=value # overrides property 'level-name'
MC_LEVEL_SEED=value # overrides property 'level-seed'
```

## Debug/Test

A command to test the configurator, on windows.

```pwsh
$env:MC_LEVEL_NAME="test world"; .\mc-config.exe; $env:MC_LEVEL_NAME=$null
```
