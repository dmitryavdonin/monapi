#---Build stage---
FROM --platform=linux/amd64 golang:1.20 AS builder
COPY . /go/src/
WORKDIR /go/src/cmd/app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o /go/bin/service

#---Final stage---
FROM --platform=linux/amd64 alpine:latest
COPY --from=builder /go/bin/service /go/bin/service
COPY --from=builder /go/src/migrations /migrations
CMD /go/bin/service --port 8000 --host '0.0.0.0'