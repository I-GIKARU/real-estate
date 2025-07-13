# Use the official Go image as build stage
FROM golang:1.23-alpine AS builder

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag CLI tool for generating docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init -g cmd/server/main.go

# Build the application with proper module settings (with Swagger enabled)
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -installsuffix cgo -o main ./cmd/server

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

# Set environment variables for Cloud Run
ENV PORT=8080
ENV SWAGGER_HOST=real-estate-backend-840370620772.us-central1.run.app

# Command to run the application
CMD ["./main"]
