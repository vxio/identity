
USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

GEN_CODE_LOCATION := "pkg/gen"

build: openapitools


# Generate the go code from the public and internal api's
openapitools:
	rm -rf pkg/gen

	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/cfg/client-generator-config.yml /local/cfg/server-generator-config.yml
