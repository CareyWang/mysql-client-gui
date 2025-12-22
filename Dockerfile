FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies (make and bash)
RUN apk add --no-cache git bash make

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download || go mod download -v

# Copy source code and Makefile
COPY . .

# Build using Makefile
RUN make build

# Final stage - create minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/mysql-client-go /app/mysql-client-go

# Set entrypoint
ENTRYPOINT ["/app/mysql-client-go"]
