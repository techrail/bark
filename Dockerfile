FROM golang:1.21-bookworm

WORKDIR /app
COPY . .
RUN go build -o bark main.go
CMD ["/app/bark"]
