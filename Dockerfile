# syntax=docker/dockerfile:1

FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o commander .

FROM alpine

WORKDIR /build

COPY .env .

COPY --from=builder /build/commander /build/commander

EXPOSE 8080

CMD [ "./commander" ]