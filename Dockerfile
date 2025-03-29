FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN apk --no-cache add \
    libc-dev gcc bash git \
    && rm -rf /var/cache/apk/*

ADD go.mod go.sum /app/
RUN go mod download

ADD . /app

RUN /app/build.sh 

FROM alpine

EXPOSE 8080

RUN apk --no-cache add tzdata curl \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist /app/bin

ENTRYPOINT ["/app/bin/api"]