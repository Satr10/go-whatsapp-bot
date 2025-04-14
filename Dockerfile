# ---------------------------
# Stage 1: Builder with CGO enabled
# ---------------------------
FROM golang:1.24-alpine AS builder

# Install gcc, musl-dev, and sqlite-dev so go-sqlite3 can compile properly
RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory for building the app
WORKDIR /app

# Copy dependency declarations and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application
COPY . .

# Build the binary with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o whatsappbot ./cmd/whatsappbot/

# ---------------------------
# Stage 2: Final Runtime Image
# ---------------------------
FROM alpine:latest

# Install CA certificates and the SQLite runtime libraries
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/whatsappbot .

# Start the app
ENTRYPOINT ["./whatsappbot"]

