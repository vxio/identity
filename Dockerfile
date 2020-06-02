FROM golang:1.14-buster as builder
WORKDIR /go/src/github.com/moov-io/identity
RUN apt-get update && apt-get install make gcc g++
COPY . .
RUN go mod download
RUN make install
RUN make identity

FROM debian:10
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /go/src/github.com/moov-io/identity/bin/identity /bin/identity

ENV SQLITE_DB_PATH /data/paygate.db
# RUN adduser -q --gecos '' --disabled-login --shell /bin/false moov
# RUN chown -R moov: /data
# USER moov

EXPOSE 8200/tcp
EXPOSE 8201/tcp

VOLUME [ "/data", "/configs" ]

ENTRYPOINT ["/bin/identity"]
