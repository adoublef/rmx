# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
WORKDIR /usr/src

ARG EXE_NAME
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

COPY go.* .
RUN go mod download

COPY . .
RUN go build \
    -ldflags "-s -w -extldflags '-static'" \
    -buildvcs=false \
    -o /usr/local/bin/ ./...

FROM alpine:${ALPINE_VERSION} AS runtime
WORKDIR /opt


COPY --from=build /usr/local/bin/rmx ./a

# run migrations
# note - use an env variable somehow
RUN ./a migrate --dsn ./rehearsal.db

CMD ["./a", "serve", "--dsn" , "./rehearsal.db"]