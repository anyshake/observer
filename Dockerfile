FROM alpine:latest AS web
# Uncomment the following line to use a mirror of APK repository
# ENV APK_SOURCE_HOST="mirror.bjtu.edu.cn"
# Uncomment the following line to use a mirror of npm registry
# ENV NPM_REGISTRY_HOST="registry.npmmirror.com"
COPY . /build_src
WORKDIR /build_src/web/src
RUN if [ ! -d "../dist" ]; then \
    if [ "x${APK_SOURCE_HOST}" != "x" ]; then \
    sed -i "s/dl-cdn.alpinelinux.org/$APK_SOURCE_HOST/g" /etc/apk/repositories; \
    fi \
    && apk add --update --no-cache nodejs npm \
    && if [ "x${NPM_REGISTRY_HOST}" != "x" ]; then \
    npm config set registry https://$NPM_REGISTRY_HOST; \
    fi \
    && npm config set loglevel=http \
    && npm ci \
    && npm run build; \
    fi

FROM golang:alpine AS builder
# Uncomment the following line to use a mirror of APK repository
# ENV APK_SOURCE_HOST="mirror.bjtu.edu.cn"
# Uncomment the following line to use a mirror of go module proxy
# ENV GOPROXY="https://goproxy.cn,direct"
COPY . /build_src
COPY --from=web /build_src/web/dist /build_src/web/dist
WORKDIR /build_src
RUN if [ "x${APK_SOURCE_HOST}" != "x" ]; then \
    sed -i "s/dl-cdn.alpinelinux.org/$APK_SOURCE_HOST/g" /etc/apk/repositories; \
    fi \
    && apk add --update --no-cache git make \
    && BUILD_PLATFORM=dockerfile make build

FROM scratch
COPY --from=builder /build_src/build/dist/observer /observer
ENTRYPOINT [ "/observer" ]
