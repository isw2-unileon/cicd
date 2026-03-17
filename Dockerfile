# ==============================================================================
# Dockerfile — multi-stage build for the calculator API
#
# Teaching note: Multi-stage builds are best practice for Go apps.
# Stage 1 (builder): full Go toolchain, compiles the binary.
# Stage 2 (final):   minimal runtime image, only the compiled binary.
# Result: a tiny, secure production image (~10MB vs ~800MB).
# ==============================================================================

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

# Install ca-certificates so the binary can make HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy dependency files first — Docker caches this layer if go.mod/go.sum
# haven't changed, so repeated builds skip the `go mod download` step.
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# ── Stage 2: Run ──────────────────────────────────────────────────────────────
FROM scratch

# Copy CA certs so HTTPS works inside the container
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the compiled binary — no Go toolchain, no source code
COPY --from=builder /bin/server /bin/server

EXPOSE 8080

ENTRYPOINT ["/bin/server"]
