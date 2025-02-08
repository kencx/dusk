FROM golang:1.23-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd && \
    go build -ldflags "-s -w -X main.version=test" -tags "linux libsqlite3 fts5" -o dusk .

FROM alpine AS final
WORKDIR /app

RUN apk add --no-cache sqlite-dev
COPY --from=builder /app/cmd/dusk /app/dusk
ENTRYPOINT ["/app/dusk"]
