FROM golang:1.21-bookworm

WORKDIR /app
COPY . .
RUN go build -o bark main.go
# You would have to still use `-p 8080:8080` switch when doing `run` or `exec`
EXPOSE 8080
CMD ["/app/bark"]
