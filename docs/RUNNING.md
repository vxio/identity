# Identity
**[Purpose](README.md)** | **[Configuration](CONFIGURATION.md)** | **Running** | **[Client](../pkg/client/README.md)**

## Running

### Getting Started

As this is only a part of a fully auth system it requires other pieces to be ran as well. Visit to our [Auth Example](https://github.com/moov-io/auth-example) project to get Identity, AuthN, and Tumbler running.

Identity can be ran various ways via docker or running locally. However you will have to run a service like `authn` to generate a token that you will use to hit `/authenticated` or `/register` the user. You will also need a service like `Tumbler` that generates the Gateway token thats verified on every other endpoint. Its best to follow our [Auth Example](https://github.com/moov-io/auth-example) project.

More tutorials to come on how to use this as other pieces required to handle authorization are in place!

- [Using docker-compose](#local-development)
- [Using our Docker image](#docker-image)

No configuration is required to serve on `:8200` and metrics at `:8201/metrics` in Prometheus format.

### Docker image

You can download [our docker image `moov/identity`](https://hub.docker.com/r/moov/identity/) from Docker Hub or use this repository. 

### Local Development

```
make run
```

---
**[Next - Client](../pkg/client/README.md)**