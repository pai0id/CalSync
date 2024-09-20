FROM golang:1.23.1-alpine

RUN apk add --no-cache make

RUN mkdir /build
WORKDIR /build

COPY ./go.* .
RUN go mod download

COPY ./cmd/. cmd/
COPY ./Makefile .
COPY ./internal/. internal/
COPY ./env/. env/

CMD ["make", "run"]