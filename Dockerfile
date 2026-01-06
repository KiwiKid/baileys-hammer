# Start from a specific golang base image
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container and build the Go app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Start a new stage from a specific debian slim version for consistency
FROM golang:1.23.4

# Install ca-certificates in one layer to reduce size
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# ==== litestream ====

# Install wget to download Litestream
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

# Download and install Litestream
ARG LITESTREAM_VERSION=v0.3.9
RUN wget https://github.com/benbjohnson/litestream/releases/download/${LITESTREAM_VERSION}/litestream-${LITESTREAM_VERSION}-linux-amd64.tar.gz && \
    tar -C /usr/local/bin -xzf litestream-${LITESTREAM_VERSION}-linux-amd64.tar.gz && \
    rm litestream-${LITESTREAM_VERSION}-linux-amd64.tar.gz


COPY litestream.yml /etc/litestream.yml
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY --from=builder /app/main /usr/local/bin/

# Command to run the executable
ENTRYPOINT ["/entrypoint.sh"]