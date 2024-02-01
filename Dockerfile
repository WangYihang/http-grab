# syntax=docker/dockerfile:1

FROM golang:1.21.6-alpine3.19
RUN go install github.com/WangYihang/http-grab@e76f595
ENTRYPOINT [ "http-grab" ]