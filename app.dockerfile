# Environment
# todo
FROM golang:1.17 as build-env

RUN mkdir -p /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build-app

# Release
FROM alpine:latest

COPY --from=build-env /src/bin/app /bin/app
ENTRYPOINT ["/bin/app"]