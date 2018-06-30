FROM golang:1.10 as builder

WORKDIR /go/src/github.com/qichengzx/seqsvr

COPY . .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o seqsvr .



FROM alpine:latest

WORKDIR /go/src/seqsvr

COPY --from=builder /go/src/github.com/qichengzx/seqsvr/seqsvr ./
COPY --from=builder /go/src/github.com/qichengzx/seqsvr/config.yml ./

RUN apk add --no-cache tzdata

ENV TZ "Asia/Shanghai"

EXPOSE 8000

ENTRYPOINT ["./seqsvr"]
