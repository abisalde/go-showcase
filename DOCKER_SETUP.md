# gRPC Showcase with Docker, Traefik, and TLS

A production-ready gRPC service setup with Docker Compose, Traefik reverse proxy, TLS encryption, and hot reload using Air.

## Features

- ✅ **gRPC Server & Client** - Full gRPC implementation with Church service
- ✅ **Docker Compose** - Multi-container orchestration
- ✅ **Traefik** - Reverse proxy and load balancer
- ✅ **TLS/SSL** - Secure communication with self-signed certificates
- ✅ **Hot Reload** - Air-verse for automatic reload on code changes
- ✅ **Interceptors** - Request logging and header validation

## Prerequisites

- Docker and Docker Compose
- Go 1.25+
- Make (optional, for convenience commands)
- OpenSSL (for certificate generation)

## Quick Start

### 1. Generate SSL Certificates

```bash
make ssl
# or
./scripts/generate-ssl.sh
```

### 2. Generate Protobuf Code

```bash
make proto
```

### 3. Start Services

```bash
make up
# or
docker-compose up -d
```

### 4. View Logs

```bash
# All services
make logs

# Server only
make server-logs

# Client only
make client-logs
```

## Available Make Commands

```bash
make help          # Show all available commands
make ssl           # Generate SSL certificates
make proto         # Generate protobuf code
make build         # Build Docker images
make up            # Start all services
make down          # Stop all services
make logs          # Show logs from all services
make server-logs   # Show server logs only
make client-logs   # Show client logs only
make restart       # Restart all services
make clean         # Clean up everything
make dev           # Start full development environment
```

## Project Structure

```
.
├── client/                 # gRPC client
│   ├── main.go
│   ├── create.go
│   ├── header_interceptor.go
│   └── log_interceptor.go
├── server/
│   ├── cmd/               # Server main application
│   │   ├── main.go
│   │   └── create.go
│   ├── pkg/               # Server packages
│   │   └── middleware/
│   └── proto/             # Protocol buffer definitions
│       ├── church.proto
│       ├── church.pb.go
│       └── church_grpc.pb.go
├── ssl/                   # SSL certificates
│   ├── ca.crt
│   ├── ca.key
│   ├── server.crt
│   ├── server.key
│   └── server.pem
├── scripts/
│   └── generate-ssl.sh    # SSL generation script
├── .air.server.toml       # Air config for server
├── .air.client.toml       # Air config for client
├── Dockerfile.server      # Server Dockerfile
├── Dockerfile.client      # Client Dockerfile
├── docker-compose.yml     # Docker Compose configuration
└── Makefile              # Build commands
```

## Services

### gRPC Server

- **Port**: 50052
- **Features**:
  - TLS encryption
  - Request logging interceptor
  - Header validation interceptor
  - Hot reload with Air

### gRPC Client

- **Connected to**: grpc-server:50052
- **Features**:
  - TLS certificate validation
  - Hot reload with Air
  - Interactive mode

### Traefik

- **Dashboard**: http://localhost:8080
- **gRPC Port**: 50052
- **Features**:
  - Automatic service discovery
  - Load balancing
  - HTTP/2 support for gRPC

## Development

### Hot Reload

Both server and client use Air for hot reload. Any changes to Go files will automatically trigger a rebuild.

**Server** watches: `server/` directory  
**Client** watches: `client/` directory

### Environment Variables

**Server:**

- `GRPC_PORT` - gRPC server port (default: 50052)
- `TLS_ENABLED` - Enable/disable TLS (default: true)

**Client:**

- `GRPC_SERVER_ADDRESS` - Server address (default: localhost:50052)
- `TLS_ENABLED` - Enable/disable TLS (default: true)

## gRPC Service

### ChurchService

**Methods:**

- `CreateChurch` - Create a new church
- `GetChurch` - Get church by ID
- `ListChurches` - List all churches (paginated)
- `UpdateChurch` - Update church details
- `DeleteChurch` - Delete a church

### Example Usage

The client will automatically connect to the server on startup. Check the logs to see the interaction.

## TLS/SSL

Self-signed certificates are used for development. The setup includes:

- **CA Certificate** (`ca.crt`) - Trust certificate for clients
- **Server Certificate** (`server.crt`) - Server identity
- **Server Private Key** (`server.pem`) - Server encryption key

### Certificate Files:

- **Share**: `ca.crt` (with clients)
- **Keep Private**: `ca.key`, `server.key`, `server.pem`, `server.crt`

## Troubleshooting

### Port Already in Use

```bash
# Find and kill process using port 50052
lsof -ti:50052 | xargs kill -9
```

### Certificate Errors

```bash
# Regenerate certificates
make ssl
make restart
```

### Container Issues

```bash
# Clean everything and restart
make clean
make dev
```

### View Container Status

```bash
docker-compose ps
```

## Production Considerations

For production deployment:

1. **Replace self-signed certificates** with certificates from a trusted CA
2. **Add authentication** - Implement proper auth middleware
3. **Rate limiting** - Configure Traefik rate limits
4. **Monitoring** - Add Prometheus/Grafana
5. **Secrets management** - Use Docker secrets or vault
6. **Resource limits** - Set memory/CPU limits in docker-compose
7. **Health checks** - Implement gRPC health checking

## License

MIT
