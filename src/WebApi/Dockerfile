# Build step
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

# Cache go mod downloads and build dependencies first.
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations
RUN go build -ldflags="-s -w" -o rinha2-back-end-go .

# Final stage: use a minimal image
FROM alpine:3.21

WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/rinha2-back-end-go .

# Expose port 8080 as defined in docker-compose.yml
EXPOSE 8080

ENTRYPOINT ["./rinha2-back-end-go"]