# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS baseline
WORKDIR /usr/src

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# required for go-sqlite3 & vitejs
RUN apk update \ 
    && apk add --no-cache --upgrade --latest gcc musl-dev \
    && apk add nodejs npm

COPY go.* .
RUN go mod download

COPY . .

FROM baseline AS testing

RUN go test -v -count=1 ./...

FROM baseline AS build

RUN go generate ./... && go build \
    -ldflags "-s -w" \
    -buildvcs=false \
    -o /usr/local/bin/ ./...

FROM alpine:${ALPINE_VERSION} AS runtime
WORKDIR /opt

ARG GID=10001
ARG UID=10001
ARG GROUP=rmx
ARG USER=rmx
ARG SQL_DIR=/data/sqlite

RUN addgroup -g ${GID} ${GROUP} \
    && adduser -G ${GROUP} -u ${UID} ${USER} -D \
    && mkdir -p ${SQL_DIR} \
    && chown -R ${GROUP}:${USER} ${SQL_DIR}

COPY --from=build /usr/local/bin/rmx ./a

USER rmx

ENTRYPOINT ["./a"]
CMD ["s"]

