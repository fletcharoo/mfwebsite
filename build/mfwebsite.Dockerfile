FROM golang:latest AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd/mfwebsite ./
RUN go build main.go

FROM debian:latest
COPY --from=build /app/main /app/main
ENTRYPOINT ["/app/main"]
