FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tcp-forwarder .

RUN chmod +x tcp-forwarder

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/tcp-forwarder .
CMD ["./tcp-forwarder"]
