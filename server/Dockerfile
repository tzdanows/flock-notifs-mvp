# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.23.3 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Cache dependencies by copying go.mod and go.sum first, and downloading modules
COPY go.mod go.sum ./

# Install the Go dependencies
RUN go mod download

# Copy the entire project into the container's working directory
COPY . .

# Build the "server" executable
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o server ./server/main.go

# Stage 2: Create a lightweight runtime image
FROM alpine:3.14

# Copy the built executable from the builder stage
COPY --from=builder /app/server /app/server

WORKDIR /app

# Expose any necessary ports (e.g., 8080 for an HTTP API)
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/app/server"]