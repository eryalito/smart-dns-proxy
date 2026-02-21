# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o smart-dns-proxy .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests to external data source
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /build/smart-dns-proxy .

# Expose DNS port
EXPOSE 53/udp

# Run the application
ENTRYPOINT ["./smart-dns-proxy"]
CMD ["-host", "0.0.0.0", "-port", "53"]
