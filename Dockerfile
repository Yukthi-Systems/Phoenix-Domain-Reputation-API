FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
   -trimpath \
   -ldflags="-s -w" \
   -o /app/bin/server \
   ./cmd/server

# ---- Final stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates 

WORKDIR /app

COPY --from=builder /app/bin/server ./server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD pgrep -f "./server" || exit 1

ENTRYPOINT ["./server"]