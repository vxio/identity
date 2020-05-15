FROM golang:1.14-buster AS build
WORKDIR /build
RUN apt-get update && apt-get install make gcc g++
COPY . .
RUN make install
RUN make identity

FROM debian:10 AS runtime
WORKDIR /
COPY --from=build /build/bin/identity /app/identity

VOLUME [ "/data", "/configs" ]

#COPY configs/*jwks* /configs/

EXPOSE 8200/tcp
EXPOSE 8201/tcp
ENTRYPOINT ["/app/identity"]