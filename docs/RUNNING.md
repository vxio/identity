# Identity
**[Purpose](README.md)** | **[Configuration](CONFIGURATION.md)** | **Running** | **[Client](../pkg/client/README.md)**

--- 

## Running

### Getting Started

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