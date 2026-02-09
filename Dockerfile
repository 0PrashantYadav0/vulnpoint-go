# Stage 1: Build Stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o vulnpilot ./cmd/server/main.go

# Stage 2: Runtime Stage
FROM alpine:latest

# Install runtime dependencies and security tools
RUN apk add --no-cache \
    ca-certificates \
    nmap \
    nmap-scripts \
    nikto \
    git \
    python3 \
    py3-pip \
    curl \
    wget \
    ruby \
    ruby-dev \
    libffi-dev \
    build-base \
    && gem install wpscan --no-document \
    && apk del ruby-dev libffi-dev build-base \
    && rm -rf /var/cache/apk/*

# Install SQLMap
RUN git clone --depth 1 https://github.com/sqlmapproject/sqlmap.git /opt/sqlmap \
    && ln -s /opt/sqlmap/sqlmap.py /usr/local/bin/sqlmap

# Install Gobuster
RUN wget https://github.com/OJ/gobuster/releases/download/v3.6.0/gobuster_Linux_x86_64.tar.gz \
    && tar -xzf gobuster_Linux_x86_64.tar.gz \
    && mv gobuster /usr/local/bin/ \
    && rm gobuster_Linux_x86_64.tar.gz

# Create non-root user
RUN addgroup -g 1000 vulnpilot && \
    adduser -D -u 1000 -G vulnpilot vulnpilot

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/vulnpilot .

# Create necessary directories
RUN mkdir -p logs scan_results data/embeddings data/analysis && \
    chown -R vulnpilot:vulnpilot /app

# Switch to non-root user
USER vulnpilot

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./vulnpilot"]