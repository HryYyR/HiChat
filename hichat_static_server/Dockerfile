FROM golang:alpine
ENV GOPROXY https://goproxy.cn,direct

WORKDIR  /app
COPY go.mod ./
COPY go.sum ./


RUN go mod download

COPY . .

RUN go build -o hichat_static_server .

EXPOSE 3005

CMD ["./hichat_static_server", "-d", "docker"]