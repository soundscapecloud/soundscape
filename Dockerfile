FROM alpine:latest
MAINTAINER Streamlist <github@streamlist.cloud>

RUN apk --no-cache add \
    curl \
    ffmpeg \
    wget

WORKDIR /data

COPY streamlist-linux-amd64 /usr/bin/streamlist

ENTRYPOINT ["/usr/bin/streamlist"]
