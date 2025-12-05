# Multi-stage Dockerfile for Memos with Project Management
# Stage 1: Build Frontend
FROM node:18-alpine AS frontend
WORKDIR /workspace

# Install pnpm
RUN npm install -g pnpm@8.15.0

# Copy the web directory
COPY web/ ./web/

# Install dependencies and build
WORKDIR /workspace/web
RUN pnpm install --frozen-lockfile

# Build frontend (outputs to ../server/router/frontend/dist)
RUN pnpm run release

# Stage 2: Build Backend with embedded frontend
FROM golang:1.21-alpine AS backend
WORKDIR /backend-build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend assets from frontend stage
COPY --from=frontend /workspace/server/router/frontend/dist ./server/router/frontend/dist

# Build the backend
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o memos ./cmd/memos

# Stage 3: Final runtime image
FROM alpine:latest AS runtime
WORKDIR /usr/local/memos

# Install runtime dependencies
RUN apk add --no-cache tzdata ca-certificates

# Set timezone
ENV TZ="UTC"

# Copy binary from backend stage
COPY --from=backend /backend-build/memos /usr/local/memos/

# Create data directory
RUN mkdir -p /var/opt/memos
VOLUME /var/opt/memos

# Expose port
EXPOSE 5230

# Set environment variables
ENV MEMOS_MODE="prod"
ENV MEMOS_PORT="5230"

# Run as non-root user for security
RUN addgroup -g 1000 memos && \
    adduser -D -u 1000 -G memos memos && \
    chown -R memos:memos /usr/local/memos /var/opt/memos

USER memos

ENTRYPOINT ["./memos"]
