FROM golang:1.20-alpine as builder

ENV CGO_ENABLED 0
ENV GO111MODULE=on

WORKDIR /app

RUN apk update --no-cache && apk add --no-cache tzdata

COPY . ./

RUN go mod download

RUN go build -o main .

FROM alpine

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /app/main /app/main

RUN apk update --no-cache && apk add --no-cache ca-certificates
RUN set -ex; \
    apk add --no-cache supervisor && \
    echo "Asia/Shanghai" > /etc/timezone && \
    rm -rf /var/cache/apk/*

EXPOSE 8080
ENTRYPOINT ["/app/main"]
