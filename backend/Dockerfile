# ==========================================
# Stage 1: Build the Go binary
# ==========================================
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools (git for fetching dependencies)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker layer caching
# Note: We haven't created go.sum yet, so we use a wildcard to avoid errors if it's missing
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the actual source code
COPY . .

# Build the Go binary.
# CGO_ENABLED=0 ensures a static binary.
# -ldflags="-w -s" strips debugging info to make the binary smaller.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/insta-gate-api ./cmd/api

# ==========================================
# Stage 2: Create the minimal production image
# ==========================================
FROM alpine:3.19

WORKDIR /app

# Install root certificates for HTTPS calls and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user for security
RUN adduser -D -g '' appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /app/insta-gate-api .

# Change ownership of the binary to the non-root user
RUN chown appuser:appuser /app/insta-gate-api

# Switch to the non-root user
USER appuser

# Expose the port (Fly.io will map this)
EXPOSE 8080

# Command to run the executable
CMD ["./insta-gate-api"]
