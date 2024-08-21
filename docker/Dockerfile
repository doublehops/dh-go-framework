FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN #go mod vendor

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/run.go

EXPOSE 8089

# Run
CMD ["./main"]