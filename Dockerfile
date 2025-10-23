# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o main .

# Stage 2: Create the final, minimal container
FROM alpine:3.18

WORKDIR /root/

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# --- NEW: Copy the frontend assets ---
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose the port the app runs on
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./main"]