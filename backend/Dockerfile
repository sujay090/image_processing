# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    vips-dev \
    pkgconfig

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries
RUN CGO_ENABLED=1 go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=1 go build -o /bin/worker ./cmd/worker

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache \
    vips \
    ca-certificates

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/worker /bin/worker
COPY migrations/ /migrations/

EXPOSE 8080

CMD ["/bin/api"]
