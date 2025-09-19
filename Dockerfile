# Multi-stage Dockerfile for DocGen Service
# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Use CGO_ENABLED=0 for static binary compatible with distroless
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docgen-server ./cmd/server

# Stage 2: Create minimal runtime image
FROM gcr.io/distroless/static:nonroot

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/docgen-server .

# Copy required assets
COPY --from=builder /app/assets ./assets

# Use nonroot user for security
USER nonroot:nonroot

# Set default environment variables
ENV PORT=8080
ENV DOCGEN_SHELL_PATH=./assets/shell/template_shell.docx
ENV DOCGEN_COMPONENTS_DIR=./assets/components/
ENV DOCGEN_SCHEMA_PATH=./assets/schemas/rules.cue

# Expose the port
EXPOSE 8080

# Set the entrypoint to run in server mode
ENTRYPOINT ["./docgen-server", "-server"]