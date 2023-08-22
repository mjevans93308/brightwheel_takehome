FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o webserver cmd/cmd.go

FROM alpine:latest

COPY --from=builder /app/webserver /usr/local/bin/webserver

EXPOSE 8080
CMD ["webserver", "--port", "8080"]