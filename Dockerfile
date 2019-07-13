# Multi-stage build
FROM golang:alpine AS builder
WORKDIR /dns-lookup-tool
ADD . /dns-lookup-tool
RUN apk update && apk add git && apk add ca-certificates
RUN cd /dns-lookup-tool && go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o dns-lookup-tool

# Second stage, smaller image
FROM alpine
ENV HTTP_PORT=8080
WORKDIR /dns-lookup-tool
COPY --from=builder /dns-lookup-tool/ /dns-lookup-tool/
ENTRYPOINT ./dns-lookup-tool
