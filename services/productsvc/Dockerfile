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
