FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY cmd/server/ ./cmd/server
COPY internal/pkg/messages/ ./internal/pkg/messages
COPY internal/pkg/server/ ./internal/pkg/server

RUN GOOS=linux go build -o /moon-server ./cmd/server/main.go

EXPOSE 8080
EXPOSE 4040

CMD ["/moon-server"]
