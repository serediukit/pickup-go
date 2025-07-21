# pickup-srv

A gRPC service for picking up users from PostgreSQL database.

## Features

- gRPC service with protobuf definitions
- PostgreSQL database integration
- User filtering by age, search age, gender, search gender, location
- Docker and Docker Compose support
- Sample data initialization

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)

## Quick Start

1. **Clone and setup the project:**
```bash
git clone <repository-url>
cd pickup-srv