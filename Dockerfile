# Stage 1: Build the Go executable
FROM golang:1.25-alpine AS builder

# Install git if your project has dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the executable
RUN go build -o myapp ./cmd

# Stage 2: Create a minimal image with the executable
FROM alpine:3.20

# Create a directory for the executable
WORKDIR /app

# Copy the executable from the builder
COPY --from=builder /app/myapp .

# Expose port if your app is a server
EXPOSE 9000

# Default command
CMD ["./myapp"]