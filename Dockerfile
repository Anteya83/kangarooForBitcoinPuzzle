FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kangaroo main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/kangaroo .

ENTRYPOINT ["./kangaroo"]