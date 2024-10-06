# Stage 1: Build the Go application
FROM golang:1.22.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY api-gateway/go.mod api-gateway/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY api-gateway .

# Build the Go app
RUN go build -o api_gateway .

# Stage 2: Create a small image with only the necessary runtime dependencies
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api_gateway .

# Copy the RSA public key
COPY api-gateway/public_key.pem /app/public_key.pem

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./api_gateway"]