
FROM golang:1.21.3

WORKDIR /app

ADD . /app

RUN go mod tidy && \
    go mod download && \
    go build -o notifier .

EXPOSE 8080

CMD ["./notifier"]
