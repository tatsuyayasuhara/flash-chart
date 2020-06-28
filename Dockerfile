FROM golang:1.13-alpine3.12

RUN apk --update add --no-cache git && \
    go get -u github.com/oxequa/realize

WORKDIR /go/src/flash-chart

CMD ["realize", "start"]
