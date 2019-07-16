# Multi-stage build
# Build production vue.js app
FROM node:lts-alpine AS frontend-builder
WORKDIR /dns-lookup-tool
ADD . /dns-lookup-tool
RUN npm --prefix frontend install
RUN npm run --prefix frontend build

FROM golang:alpine AS builder
WORKDIR /dns-lookup-tool
COPY --from=frontend-builder /dns-lookup-tool/ /dns-lookup-tool/
RUN apk update && apk add git && apk add ca-certificates
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o dns-lookup-tool

# Copy final production vue files and Go binary to minimal container
FROM alpine
WORKDIR /dns-lookup-tool
COPY --from=builder /dns-lookup-tool/frontend/dist /dns-lookup-tool/frontend/dist/
COPY --from=builder /dns-lookup-tool/dns-lookup-tool /dns-lookup-tool/bin/

ENTRYPOINT ./bin/dns-lookup-tool
