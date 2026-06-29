# syntax=docker/dockerfile:1.7

FROM golang:1.26.3-alpine AS builder

WORKDIR /src

# Cache dependencies first
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
WORKDIR /src/apps/api
RUN go mod download

# Copy API source and build
COPY apps/api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/lotto-api app/main.go

# Install golang-migrate
RUN CGO_ENABLED=0 GOOS=linux GOBIN=/out go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2

FROM alpine:3.22 AS runner

RUN addgroup -S app && adduser -S app -G app && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /out/lotto-api /app/lotto-api
COPY --from=builder /out/migrate /app/migrate
COPY apps/api/migrations /app/migrations

# Non-secret runtime defaults (Fly can override)
ENV APP_ENV=production
ENV PORT=:8080

EXPOSE 8080

USER app
ENTRYPOINT ["/app/lotto-api"]
