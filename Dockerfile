# Build stage
FROM golang:1.24-alpine AS builder

# Install protoc and protoc-gen-go
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install Make
RUN apk add --no-cache make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf files
RUN make proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pickup-srv ./cmd/server/grpc_server.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/pickup-srv .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./pickup-srv"]