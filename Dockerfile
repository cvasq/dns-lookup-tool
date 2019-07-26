# Multi-stage build 

# Build production vue.js app
FROM node:lts-alpine AS frontend-builder
ARG VUE_APP_BASE_PATH
ENV VUE_APP_BASE_PATH $VUE_APP_BASE_PATH
WORKDIR /dns-lookup-tool
ADD . /dns-lookup-tool
RUN npm --prefix ui install
RUN npm run --prefix ui build

# Build go binary
FROM golang:alpine AS builder
ENV GO111MODULE=on
WORKDIR /dns-lookup-tool
COPY --from=frontend-builder /dns-lookup-tool/ /dns-lookup-tool/
RUN apk update && apk add git && apk add ca-certificates
RUN go get -d -v
RUN go get github.com/rakyll/statik
RUN statik -src=./ui/dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o dns-lookup-tool

# Copy final build to minimal container
FROM alpine
WORKDIR /dns-lookup-tool
COPY --from=builder /dns-lookup-tool/dns-lookup-tool ./dns-lookup-tool
ENTRYPOINT ./vue-websocket-echo
