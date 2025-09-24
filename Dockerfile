FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY day03/party/ ./day03/party/
WORKDIR /app/day03/party
RUN go build -o /app/app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/day03/party/*.html .
EXPOSE 5000
CMD ["./app"]
