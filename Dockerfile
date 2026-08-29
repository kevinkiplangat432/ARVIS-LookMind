FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o arvis ./cmd/arvis/

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/arvis .
COPY --from=builder /app/migrations ./migrations

# Copy the entrypoint script
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8080 8081

ENTRYPOINT ["/app/docker-entrypoint.sh"]