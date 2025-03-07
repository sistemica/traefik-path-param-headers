FROM golang:1.20-alpine

WORKDIR /app

COPY server.go .

RUN go mod init backend && \
    go mod tidy && \
    go build -o server .

EXPOSE 8000

CMD ["./server"]