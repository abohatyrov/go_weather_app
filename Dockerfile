# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21.5
ARG ALPINE_VERSION=3.18
ARG OPENWEATHERMAP_API_KEY

### Base image
##################################################
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base

WORKDIR /usr/src/app
RUN apk add --no-cache git xdg-utils && \
    apk update && \
    apk upgrade

### Install dependencies
##################################################
FROM base AS deps

COPY go.mod .
COPY go.sum .

RUN go mod download

### Build
##################################################
FROM deps AS build

COPY . .
RUN go build -o /usr/src/app

### Final image
##################################################
FROM alpine:${ALPINE_VERSION} AS final

WORKDIR /usr/src/app

ENV OPENWEATHERMAP_API_KEY=${OPENWEATHERMAP_API_KEY}

COPY --from=build /usr/src/app/templates ./templates
COPY --from=build /usr/src/app/static ./static
COPY --from=build /usr/src/app/openweathermap .

EXPOSE 8080

CMD ["./openweathermap"]