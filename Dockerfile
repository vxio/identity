FROM debian:buster AS runtime
WORKDIR /

RUN apt-get update && apt-get install -y ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY bin/.docker/identity /app/identity
VOLUME [ "/data", "/configs" ]

EXPOSE 8200/tcp
EXPOSE 8201/tcp

VOLUME [ "/data", "/configs" ]

ENTRYPOINT ["/app/identity"]
