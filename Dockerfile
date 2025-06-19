# ---- Build Stage ----
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Cache go.mod & go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o labbi-app ./cmd

# ---- Final Stage ----
FROM alpine:latest

WORKDIR /app

# Install CA certificates for TLS
RUN apk add --no-cache ca-certificates

# Copy the binary
COPY --from=builder /app/labbi-app ./labbi-app

# Copy static assets and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/internal/templates ./templates

# Expose port
EXPOSE 8080

# Default environment (can be overridden)
ENV SERVER_ADDRESS=":8080"

# Entry point
ENTRYPOINT ["/app/labbi-app"]
