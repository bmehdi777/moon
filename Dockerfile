# build
FROM golang:1.25 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY cmd/server/ ./cmd/server
COPY internal/pkg/communication/ ./internal/pkg/communication
COPY internal/pkg/server/ ./internal/pkg/server

RUN GOOS=linux go build -o ./server ./cmd/server/main.go

# release
FROM debian:stable-slim AS release-stage

WORKDIR /

COPY --from=build-stage /app/server /server

EXPOSE 8080
EXPOSE 4040

CMD ["/server"]
