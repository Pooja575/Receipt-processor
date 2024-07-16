# Use an official Go runtime as a parent image
FROM golang:1.17-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to ensure dependencies are downloaded efficiently
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
