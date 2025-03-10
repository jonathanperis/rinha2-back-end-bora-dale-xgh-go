# Build step
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Cache go mod downloads and build dependencies first.
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations
RUN go build -ldflags="-s -w" -o rinha-go .

# Final stage: use a minimal image
FROM alpine:latest

WORKDIR /app
# Copy the binary from the builder stage

COPY --from=builder /app/rinha-go .

# Expose port 9999 as defined in docker-compose.yml
EXPOSE 9999

ENTRYPOINT ["./rinha-go"]