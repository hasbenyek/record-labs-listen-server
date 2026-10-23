# Record Labs — Listen Together server.
# Builds the server from this repository and runs it as a small Alpine image.
# Render (or any Docker host) sets PORT; the server listens on it (default 8080).

FROM golang:1.25.6-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o recordlabs-listen ./cmd/recordlabs-listen

FROM alpine:3.20
RUN apk add --no-cache ca-certificates \
    && addgroup -g 1000 recordlabs \
    && adduser -D -u 1000 -G recordlabs recordlabs
WORKDIR /app
COPY --from=builder /build/recordlabs-listen .
COPY ua_policy.json /app/ua_policy.json
RUN mkdir /app/data && chown -R recordlabs:recordlabs /app
USER recordlabs
ENV UA_POLICY_FILE=/app/ua_policy.json
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider "http://localhost:${PORT:-8080}/health" || exit 1
CMD ["./recordlabs-listen"]
