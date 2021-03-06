FROM debian:buster AS runtime
LABEL maintainer="Moov <support@moov.io>"

WORKDIR /

RUN apt-get update && apt-get install -y ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY bin/.docker/identity /app/identity

EXPOSE 8200/tcp
EXPOSE 8201/tcp

VOLUME [ "/data", "/configs" ]

ENTRYPOINT ["/app/identity"]
