# Stage 1: Build with CGO enabled
FROM golang:1.23-alpine AS builder

# Install required build tools for CGO
RUN apk add --no-cache gcc g++ make sqlite-dev

# Enable CGO (required for go-sqlite3)
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy source code
COPY day03/party/ ./

# Download dependencies
RUN go mod download

# Build the Go binary with SQLite support
RUN go build -ldflags="-w -s" -o out .

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install runtime dependencies (SQLite libs + SSL certs)
RUN apk --no-cache add ca-certificates sqlite-libs

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy compiled binary and static assets
COPY --from=builder /app/out .
COPY --from=builder /app/*.html ./
COPY --from=builder /app/styles.css ./
COPY --from=builder /app/app.js ./

# Set permissions
RUN chown -R appuser:appuser /app
USER appuser

# Expose port 8080
EXPOSE 8080

# Run the app
CMD ["./out"]
