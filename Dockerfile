FROM golang:1.13 AS builder

WORKDIR /src
COPY . .

RUN make build && chmod +x /src/build/chatops

FROM scratch

COPY --from=builder /src/build/chatops /chatops
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENTRYPOINT ["/chatops"]
