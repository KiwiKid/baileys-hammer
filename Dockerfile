# Start from a specific golang base image
FROM golang:1.21.1 AS builder

WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container and build the Go app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Start a new stage from a specific debian slim version for consistency
FROM golang:1.21.1

# Install ca-certificates in one layer to reduce size
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]