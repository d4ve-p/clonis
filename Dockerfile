FROM golang:1.25.4-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -a -ldflags '-linkmode external -extldflags "-static"' -o clonis ./cmd/server

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for SSL and sqlite cli for debugging
RUN apk --no-cache add ca-certificates sqlite

# Copy binary and assets from the builder
COPY --from=builder /app/clonis .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

EXPOSE 8080

VOLUME ["/config"]

CMD ["./clonis"]