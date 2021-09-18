FROM golang:1.17 AS builder

WORKDIR /usr/local/go/src/fluent-bit-yc-logging

COPY . /usr/local/go/src/fluent-bit-yc-logging/

RUN go build -buildmode=c-shared -o fluent-bit-yc-logging.so .

FROM fluent/fluent-bit:1.8

COPY --from=builder /usr/local/go/src/fluent-bit-yc-logging/fluent-bit-yc-logging.so /fluent-bit/bin/fluent-bit-yc-logging.so
# COPY fluent-bit-local.conf /fluent-bit/etc/fluent-bit.conf

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf", "-e", "/fluent-bit/bin/fluent-bit-yc-logging.so"]