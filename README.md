# Mumble Status

A simple tool to get receive the status of your Mumble server:
- Version
- Users (online/max)
- Maximum bandwidth

## Features

- HTTP API
- Prometheus exporter

## Server

### Installation

#### Docker

The server can be run as a Docker container with
```sh
docker run -p3000:3000 -eMUMBLE_ADDRESS=<your-server>:64738 ghcr.io/lennart-k/mumble-status:stable
```

#### Manual build

You can also build the binary yourself
```
go build -o ./mumble-status-server ./cmd/server/
```

### Configuration

| Environment variable | Argument          | Description                                              |
| -------------------- | ----------------- | -------------------------------------------------------- |
| `MUMBLE_ADDRESS`     | `-mumble-address` | The address of the Mumble server (needs to include port) |
| `LISTEN_HOST`        | `-host`           | The host to listen on (default: 0.0.0.0)                 |
| `LISTEN_PORT`        | `-port`           | The port to listen on (default: 3000)                    |

If both arguments and environment variables are provided, the arguments take precedence
