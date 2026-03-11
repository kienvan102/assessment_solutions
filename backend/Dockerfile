# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Copy the data folder to ensure the JSON file is available
COPY --from=builder /app/data ./data

# Set required environment variables
ENV ENV=production
ENV PORT=8080

# Expose port 8080 to the outside world
EXPOSE 8080
CMD ["./server"]
