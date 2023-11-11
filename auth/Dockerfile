
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth_service .

# Start from a fresh Alpine image to create a smaller final image.
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth_service .

EXPOSE 8080

ENTRYPOINT ["./auth_service"]
                                       