eROM alpine:latest
    go build -v --compiler gc --ldflags "-extldflags -static -s -w -X main.version=${BUILD_VERSION}" -o /usr/bin/soundscape-linux-amd64
MAINTAINER Soundscape <soundscape@portal.cloud>

RUN apk --no-cache add \
    curl \
    ffmpeg \
    wget

WORKDIR /data

COPY soundscape-linux-amd64 /usr/bin/soundscape

ENTRYPOINT ["/usr/bin/soundscape"]
