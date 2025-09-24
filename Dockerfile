# Use the official Go image as base
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY day03/party/go.mod day03/party/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY day03/party/ ./

# Build the application
RUN go build -ldflags="-w -s" -o out .

# Use a minimal image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/out .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Run the application
CMD ["./out"]