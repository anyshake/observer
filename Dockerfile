FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git bash wget curl make npm
WORKDIR /build
RUN git clone --progress https://github.com/anyshake/observer.git ./observer && \
    export VERSION=`cat ./observer/VERSION` && \
    cd ./observer/frontend/src && \
    npm install && \
    make && \
    cd ../../docs && \
    make && \
    cd ../cmd && \
    go mod tidy && \
    go build -ldflags "-s -w -X main.version=$VERSION -X main.release=docker_build" -trimpath -o /tmp/observer *.go

FROM alpine

COPY --from=builder /tmp/observer /usr/bin/observer
RUN chmod 755 /usr/bin/observer && \
    mkdir -p /etc/observer

CMD ["observer", "-config=/etc/observer/config.json"]
