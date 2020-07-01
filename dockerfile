FROM golang:1.14.3 AS build

WORKDIR /go/cache

ENV GO111MODULE on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN go env -w GOSUMDB=off

# 公共包获取
ADD ./src/go.mod .
ADD ./src/go.sum .
RUN go mod download

WORKDIR /go/src/trex_online/src

COPY src .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o main

FROM alpine:3.8

RUN apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /go/src/trex_online

COPY --from=build /go/src/trex_online/src/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
