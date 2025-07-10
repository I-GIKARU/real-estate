# Use the official Go image as build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Use a minimal Alpine image for runtime
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations if they exist
COPY --from=builder /app/migrations ./migrations

# Expose the port that the app runs on
EXPOSE 8080

# Set environment variable for port (Cloud Run requirement)
ENV PORT=8080

# Command to run the application
CMD ["./main"]
