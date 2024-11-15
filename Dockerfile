FROM golang:latest AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd/mfwebsite ./
RUN go build main.go

FROM debian:latest
WORKDIR /app
COPY --from=build /app/main /app/main
COPY test/markdown /app
ENTRYPOINT ["/app/main"]
