FROM golang:1.22 AS builder

RUN go install github.com/goreleaser/goreleaser@latest
WORKDIR /app
COPY ./.git/ ./.git/
RUN git reset --hard HEAD
RUN goreleaser build --clean --id=http-grab --snapshot

FROM scratch
COPY --from=builder /app/dist/http-grab_linux_amd64_v1/http-grab /usr/local/bin/http-grab
ENTRYPOINT [ "/usr/local/bin/http-grab" ]