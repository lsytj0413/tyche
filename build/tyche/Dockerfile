# ipwhere multi-stage build Dockerfile

FROM alpine:3.7 AS build-env
ADD . /go/src/github.com/lsytj0413/tyche
WORKDIR /go/src/github.com/lsytj0413/tyche
ENV GOPATH /go
RUN export PATH=/go/bin:$PATH \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    ## install musl-dev for gcc
    && apk add --no-cache git go musl-dev make \
    && make \
    && mkdir /out \
    && cp ./bin/tyche /out/tyche \
    && rm -rf /var/cache/apk/*

FROM alpine:3.7

LABEL MAINTAINER sorenyang@foxmail.com

ENV GIN_MODE release
RUN apk add -U tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /out/tyche /usr/local/bin/tyche
EXPOSE 80
WORKDIR /usr/local/bin
CMD ["tyche"]
