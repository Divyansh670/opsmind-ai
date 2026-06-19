# =============================================
# Stage 1: Build the Go binary
# =============================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# Tidy and build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o opsmind-backend ./cmd/api/main.go

# =============================================
# Stage 2: Minimal runtime image
# =============================================
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS calls (Groq API, GitHub API)
RUN apk --no-cache add ca-certificates tzdata

# Copy only the compiled binary from builder
COPY --from=builder /app/opsmind-backend .

EXPOSE 8080

CMD ["./opsmind-backend"]