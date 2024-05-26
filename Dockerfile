FROM golang:1.22.1
WORKDIR /app
ADD . .
RUN go mod download
CMD ["/bin/bash", "-c", "go run main.go"]