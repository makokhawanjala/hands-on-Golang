FROM golang:1.23-alpine AS builder

# Add build metadata to break cache
LABEL build-version="2024-09-24-v2"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY day03/party/ ./day03/party/
WORKDIR /app/day03/party
RUN go build -o /app/app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the correctly named binary
COPY --from=builder /app/app .
COPY --from=builder /app/day03/party/*.html .

EXPOSE 5000
CMD ["./app"]
