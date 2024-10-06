# Stage 1: Build the Go application
FROM golang:1.22.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY chat-service/go.mod chat-service/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY chat-service .

# Build the Go app
RUN go build -o chat_service .

# Stage 2: Create a small image with only the necessary runtime dependencies
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/chat_service .

# Expose port 8082 to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["./chat_service"]