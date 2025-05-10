# Container definition for server
# server/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o algo-scales-server ./server

# Create a minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/algo-scales-server .

# Expose the API port
EXPOSE 8080

# Command to run the executable
CMD ["./algo-scales-server"]
