FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /urlshort ./cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /urlshort /urlshort

EXPOSE 8080

CMD ["/urlshort"]