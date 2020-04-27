# Identity
Handles identity management for our users who will be managing the system.


[Internal API](docs/internal-api.yml)
[Public API](docs/public-api.yml)


docker run --rm -u 1000:1000 -e OPENAPI_GENERATOR_VERSION=4.2.0 -v /home/jj/Documents/moov/moov-identity:/local openapitools/openapi-generator-cli list
