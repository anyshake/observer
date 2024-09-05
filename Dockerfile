FROM alpine:latest AS frontend
# Uncomment the following line to use a mirror of APT repository
# ENV APK_SOURCE_HOST="mirrors.bfsu.edu.cn"
# Uncomment the following line to use a mirror of npm registry
# ENV NPM_REGISTRY_HOST="registry.npmmirror.com"
COPY . /build_src
WORKDIR /build_src/frontend/src
RUN if [ ! -d "../dist" ]; then \
    if [ "x${APK_SOURCE_HOST}" != "x" ]; then \
    sed -i "s/dl-cdn.alpinelinux.org/$APK_SOURCE_HOST/g" /etc/apk/repositories; \
    fi \
    && apk add --update --no-cache nodejs npm \
    && if [ "x${NPM_REGISTRY_HOST}" != "x" ]; then \
    npm config set registry https://$NPM_REGISTRY_HOST; \
    fi \
    && npm config set loglevel=http \
    && npm install \
    && npm run build; \
    fi

FROM golang:alpine AS builder
# Uncomment the following line to use a mirror of go module proxy
# ENV GOPROXY="https://goproxy.cn,direct"
COPY . /build_src
COPY --from=frontend /build_src/frontend/dist /build_src/frontend/dist
WORKDIR /build_src/docs
RUN go get -v github.com/swaggo/swag/cmd/swag \
    && go install -v github.com/swaggo/swag/cmd/swag \
    && swag init -g ../cmd/main.go -d ../api,../config,../drivers/explorer,../server -o ./
WORKDIR /build_src/cmd
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(cat ../VERSION) -X main.tag=dockerbuild" \
    -v -trimpath \
    -o /tmp/observer

FROM scratch
COPY --from=builder /tmp/observer /observer
ENTRYPOINT [ "/observer" ]
