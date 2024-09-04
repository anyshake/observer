FROM node:alpine AS frontend
# Uncomment the following line to use a mirror of npm registry
# ENV NPM_REGISTRY="https://registry.npmmirror.com"
COPY . /build_src
WORKDIR /build_src/frontend/src
RUN if [ ! -d "../dist" ]; then \
    npm config set registry $NPM_REGISTRY \
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
