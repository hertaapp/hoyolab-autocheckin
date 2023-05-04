FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /go/src/
COPY . /go/src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/hoyolab-autocheckin .

FROM scratch
COPY --from=builder /go/bin/hoyolab-autocheckin /hoyolab-autocheckin
COPY ./certs/*.pem /etc/ssl/certs/
COPY ./static/*.mp4 /static/
ENTRYPOINT ["/hoyolab-autocheckin"]
