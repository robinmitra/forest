FROM golang:alpine

ENV GO111MODULE=on

# Disabled as Go modules seem to have trouble with CGO at the moment as not all vendors have opted
# into it. Build completed correct, but running unit tests fail with error saying 'gcc' is missing.
ENV CGO_ENABLED=0

RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

# Pre-emptively download Go module dependencies.
# This has the effect that Go doesn't download module dependencies everytime our code changes, but
# does download them in case the dependencies specified (in go.mod and go.sum) do actually change.
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build
