# TODO use builder and base images 

# docker/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go
RUN chmod +x ./main
 
EXPOSE 8080
# Removed due to the delay check command in docker compose
# CMD ["./main"]
