# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# The user will run go mod init separately, so we handle missing go.mod gracefully
# In a real build, we'd copy go.mod and go.sum here, but we'll just copy everything
COPY . .

# Run go build, creating the binary at /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Copy the data folder to ensure the JSON file is available
COPY --from=builder /app/data ./data

# Expose port 8080 to the outside world
EXPOSE 8080
CMD ["./server"]
