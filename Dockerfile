FROM golang:1.18-alpine AS builder
ADD . /build
RUN cd /build && go install -mod=mod ./cmd/nomad-event-collector

FROM alpine:latest
COPY --from=builder /go/bin/nomad_event_collector .
CMD ["./nomad_event_collector"]
