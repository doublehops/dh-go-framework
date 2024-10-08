FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/run.go

# Final stage: create a minimal image for the Go application
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.json /root/config.json
EXPOSE 8088
