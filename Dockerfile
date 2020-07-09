FROM golang:1.13

WORKDIR /src
COPY . .

RUN make build

FROM scratch

COPY --from=0 /src/build/chatops /chatops

ENTRYPOINT ["/chatops"]
