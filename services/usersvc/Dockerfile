FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o app .

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /src/app .
EXPOSE 8081
USER nonroot
CMD ["./app"]
