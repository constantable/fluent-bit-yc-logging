FROM golang:1.17 AS builder

WORKDIR /usr/local/go/src/fluent-bit-yc-logging

COPY . /usr/local/go/src/fluent-bit-yc-logging/

RUN pwd && ls -la
RUN make build

FROM fluent/fluent-bit:1.8

COPY --from=builder /usr/local/go/src/fluent-bit-yc-logging/fluent-bit-yc-logging.so /fluent-bit/bin/fluent-bit-yc-logging.so

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf", "-e", "/fluent-bit/bin/fluent-bit-yc-logging.so"]