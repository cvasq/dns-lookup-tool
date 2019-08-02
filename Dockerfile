# Multi-stage build 

# Build production vue.js app
FROM node:lts-alpine AS frontend-builder
ARG VUE_APP_WS_URL
ENV VUE_APP_WS_URL $VUE_APP_WS_URL
ARG VUE_APP_BASE_PATH
ENV VUE_APP_BASE_PATH $VUE_APP_BASE_PATH
WORKDIR /dns-lookup-tool
ADD . /dns-lookup-tool/
RUN npm --prefix ui install
RUN npm run --prefix ui build
