# Dockerfile
FROM alpine:latest

# This magic variable is passed by GoReleaser's buildx command
ARG TARGETPLATFORM

# Don't COPY go.sum or go.mod. 
# Just copy the binary that GoReleaser ALREADY built.
COPY ${TARGETPLATFORM}/clonis /usr/bin/clonis

ENTRYPOINT ["/usr/bin/clonis"]