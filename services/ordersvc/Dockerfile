# Use Go 1.22 instead of 1.21
FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o service .

# Final minimal image
FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8081
CMD ["./service"]