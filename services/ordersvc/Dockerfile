FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/app /app/app
CMD ["/app/app"]
