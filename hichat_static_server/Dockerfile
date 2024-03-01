# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

ENV GOPROXY https://goproxy.cn,direct

WORKDIR  /app
COPY go.mod ./
COPY go.sum ./
COPY ./ ./

RUN go mod download

RUN go build .

EXPOSE 3005

CMD [ "./hichat_static_server" ]