# Build Stage
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN go mod download
# Build a static binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o clonis ./cmd/server

# Run Stage
FROM alpine:latest
WORKDIR /root/

# Install ca-certificates for Google Drive SSL connections
RUN apk --no-cache add ca-certificates sqlite

COPY --from=builder /app/clonis .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Expose the port
EXPOSE 8080

# The volume for the database and config
VOLUME ["/config"]

CMD ["./clonis"]