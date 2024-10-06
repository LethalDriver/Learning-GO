# Stage 1: Build the Go application
FROM golang:1.22.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY user-service/go.mod user-service/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY user-service .

# Build the Go app
RUN go build -o user_service .

# Stage 2: Create a small image with only the necessary runtime dependencies
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user_service .

# Copy the RSA private key
COPY user-service/private_key.pem /app/private_key.pem

# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./user_service"]