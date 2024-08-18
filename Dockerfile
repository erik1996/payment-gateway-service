# Start with the official Golang image
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Set the working directory to where the main.go file is located
WORKDIR /app/cmd

# Build the Go app
RUN go build -o /app/main .

# Final stage
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache bash

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 (or replace with ${PORT} if necessary)
EXPOSE ${PORT}

# Command to run the executable
CMD ["/app/main"]
