FROM golang:1.18-alpine AS builder
ADD . /build
RUN cd /build && go install -mod=mod ./cmd/nomad-event-collector

FROM alpine:latest

ENV TZ=Europe/Moscow
RUN apk --no-cache add ca-certificates tzdata && cp -r -f /usr/share/zoneinfo/$TZ /etc/localtime

COPY --from=builder /go/bin/nomad-event-collector .
CMD ["./nomad-event-collector"]
