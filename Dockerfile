FROM golang:1.22-alpine3.18 as builder
WORKDIR /app
ENV CGO_ENABLED=0 GOFLAGS="-ldflags=-s -w"
COPY go.mod go.sum ./
RUN go mod download

RUN go vet -v && go build -v -o ./dusk cmd

FROM scratch
LABEL maintainer="kencx"

WORKDIR /
COPY --from=builder --chmod=+x /app/dusk ./
VOLUME data
EXPOSE 9090
CMD []
