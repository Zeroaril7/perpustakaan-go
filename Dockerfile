FROM golang:1.19-alpine3.16 AS builder

LABEL maintainer="fakhril"

WORKDIR /usr/src/app

# Update package
RUN apk add --update --no-cache --virtual .build-dev build-base git

COPY . .

RUN go build -o ./docker-go app/main.go

FROM alpine:3.16

WORKDIR /usr/src/app

RUN mkdir ./source

COPY --from=builder /usr/src/app/docker-go ./run-me

COPY --from=builder /usr/src/app/ ./source

USER root

EXPOSE 9050

CMD ["./run-me"]
