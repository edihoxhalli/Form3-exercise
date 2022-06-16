# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

RUN apk add build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY account account
COPY it it

COPY wait-for.sh wait-for.sh