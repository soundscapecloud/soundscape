FROM golang:1.9-alpine as build
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src/github.com/streamlist/streamlist
ARG STREAMLIST_VERSION=unknown
ENV GODEBUG="netdns=go http2server=0"
ENV GOPATH="/go"
RUN go get \
    github.com/jteeuwen/go-bindata/... \
    github.com/PuerkitoBio/goquery \
    github.com/armon/circbuf \
    github.com/disintegration/imaging \
    github.com/dustin/go-humanize \
    github.com/julienschmidt/httprouter \
    github.com/eduncan911/podcast \
    github.com/rylio/ytdl \
    go.uber.org/zap \
    golang.org/x/crypto/acme/autocert \
    github.com/jinzhu/gorm \
    github.com/jinzhu/gorm/dialects/sqlite \
    github.com/go-sql-driver/mysql
COPY *.go ./
COPY internal ./internal
COPY static ./static
COPY templates ./templates
RUN go fmt && \
    go vet --all && \
    go-bindata --pkg main static/... templates/...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -v --compiler gc --ldflags "-extldflags -static -s -w -X main.version=${STREAMLIST_VERSION}" -o /usr/bin/streamlist-linux-amd64
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 \
#    go build -v --compiler gc --ldflags "-extldflags -static -s -w -X main.version=${STREAMLIST_VERSION}" -o /usr/bin/streamlist-linux-armv7
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#    go build -v --compiler gc --ldflags "-extldflags -static -s -w -X main.version=${STREAMLIST_VERSION}" -o /usr/bin/streamlist-linux-arm64

FROM alpine:latest
RUN apk --no-cache add \
    curl \
    ffmpeg \
    wget
WORKDIR /data
COPY --from=build /usr/bin/streamlist-linux-amd64 /usr/bin/streamlist
EXPOSE 80
ENTRYPOINT ["/usr/bin/streamlist"]
