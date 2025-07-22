# pickup-srv

A gRPC service for picking up users from PostgreSQL database.

## Features

- gRPC service with protobuf definitions
- PostgreSQL database integration
- RabbitMQ consumer for user registration events
- User filtering by age, search age, gender, search gender, location
- Docker and Docker Compose support
- Sample data initialization
- Event-driven user creation

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)

## Quick Start

1. **Clone and setup the project:**

```bash
git clone <repository-url>
cd pickup-srv
```

2. **Start the service with Docker Compose:**

```bash 
make docker-up
```

This will:

- Start PostgreSQL database with sample data
- Start RabbitMQ with management interface
- Build and start the pickup-srv gRPC service with event consumer
- Expose the gRPC service on port 8080
- Expose RabbitMQ management UI on port 15672

3. **Test the service:**

```bash 
# Generate protobuf files (if not using Docker)
make proto
# Build and run the client to test gRPC
make build-grpc-client
# Build and run the producer to test RabbitMQ events
make build-user-registered-produced
```

## Services

### gRPC Service

The service provides one gRPC method:

#### GetUsers

```
protobuf rpc GetUsers(GetUsersRequest) returns (GetUsersResponse);
```

**Request:**

- `user_params`: Optional filtering parameters
    - `name`: Filter by name (case-insensitive partial match)
    - `email`: Filter by email (case-insensitive partial match)
    - `age`: Filter by exact age
    - `city`: Filter by city (case-insensitive partial match)
- `limit`: Maximum number of users to return

**Response:**

- `users`: Array of user objects
- `total`: Total number of users matching the filter

### RabbitMQ Consumer

The service listens to RabbitMQ queue `user.registration` for user registration events.

#### Event Format

```json
{
  "name": "John Doe",
  "age": 30,
  "city": "New York",
  "gender": "m",
  "search_gender": "f",
  "search_age_from": 20,
  "search_age_to": 25,
  "location": 19
}
```

## Management Interfaces

- **RabbitMQ Management UI**: http://localhost:15672 (guest/guest)
- **gRPC Service**: localhost:8080

## Development

### Local Development

1. **Start infrastructure:**

```bash
# Start PostgreSQL and RabbitMQ
docker-compose up postgres rabbitmq -d
```

2. **Initialize database:**

```
bash psql -h localhost -U postgres -d pickup_db -f init.sql
``` 

3. **Generate protobuf files:**

```
bash make proto
``` 

4. **Run the service:**

```
bash make run
``` 

5. **Test RabbitMQ events:**

```
bash make build-user-registered-produced
``` 

### Environment Variables

#### Database

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: password)
- `DB_NAME`: Database name (default: pickup_db)
- `DB_SSLMODE`: SSL mode (default: disable)

#### RabbitMQ

- `RABBITMQ_HOST`: RabbitMQ host (default: localhost)
- `RABBITMQ_PORT`: RabbitMQ port (default: 5672)
- `RABBITMQ_USER`: RabbitMQ user (default: guest)
- `RABBITMQ_PASSWORD`: RabbitMQ password (default: guest)

#### Service

- `GRPC_PORT`: gRPC service port (default: 8080)

### Make Commands

- `make proto`: Generate protobuf files
- `make build`: Build the server binary
- `make run`: Run the server locally
- `make build-user-registered-produced`: Build and run the publisher binary
- `make docker-build`: Build Docker image
- `make docker-up`: Start with Docker Compose
- `make docker-down`: Stop Docker Compose services

## Database Schema

```sql 
CREATE TABLE IF NOT EXISTS users
(
    id
    SERIAL
    PRIMARY
    KEY,
    name
    VARCHAR
(
    255
) NOT NULL,
    age INTEGER NOT NULL CHECK
(
    age >
    0
),
    city VARCHAR
(
    255
) NOT NULL,
    reg_dt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    gender VARCHAR
(
    1
) NOT NULL,
    search_gender VARCHAR
(
    1
) NOT NULL,
    search_age_from INTEGER NOT NULL CHECK
(
    search_age_from >
    0
),
    search_age_to INTEGER NOT NULL CHECK
(
    search_age_to >
    search_age_from
),
    location FLOAT NOT NULL
    )
``` 

## Message Queue

The service consumes messages from the `user.registration` queue. Each message should contain user registration data in
JSON format. The consumer will:

1. Validate the event data
2. Insert the new user into the PostgreSQL database
3. Acknowledge the message on success
4. Reject and requeue the message on failure

## Testing

1. **Start the services:**

```bash
make docker-up
```

1. **Test gRPC API:**

``` bash
go run ./cmd/client
```

2. **Send registration events:**

``` bash
make publisher
```

3. **Verify new users were created:**

``` bash
go run ./cmd/client
```

The publisher will send sample user registration events to RabbitMQ, which will be processed by the consumer and stored
in the database.

Perfect! Now you have a complete pickup-srv service with RabbitMQ event processing. Here's what I've added:

## New Components:

1. **RabbitMQ Integration** (`internal/messaging/rabbitmq.go`):
    - Connection handling with retry logic
    - Queue declaration and message consumption

2. **Event Consumer** (`internal/consumer/user_consumer.go`):
    - Listens to `user.registration` queue
    - Validates incoming user registration events
    - Creates users in the database
    - Handles message acknowledgment/rejection

3. **Updated Models**:
    - Added `UserRegistrationEvent` for RabbitMQ events
    - Updated `User` model to match your schema (with `reg_dt`)

4. **Enhanced Repository**:
    - Added `CreateUser` method for event-driven user creation

5. **Test Producer** (`test/test_producer_user_registered.go`):
    - Simulates user registration events
    - Helps test the RabbitMQ integration

## Key Features:

- **Event-Driven Architecture**: Users are automatically created when registration events are received
- **Graceful Shutdown**: Proper context handling for clean service shutdown
- **Error Handling**: Failed messages are requeued for retry
- **Validation**: Event data validation before database insertion
- **Docker Integration**: RabbitMQ included in docker-compose with management UI

## To test the new functionality:

1. **Start services**: `make docker-up`
2. **Send test events**: `make build-user-registered-produced`
3. **Query users**: `make build-grpc-client`
4. **Monitor RabbitMQ**: Visit http://localhost:15672 (guest/guest)

The service now handles both gRPC requests for user retrieval and RabbitMQ events for user registration!

```
