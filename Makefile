
USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

build: openapi-public

# Generate the go code from the public and internal api's
openapi-public:
	rm -rf generated/go

	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION=4.2.0 \
		-v ${PWD}:/local openapitools/openapi-generator-cli generate \
    	-i /local/docs/public-api.yml \
    	-g go \
    	-o /local/generated/go/public

	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION=4.2.0 \
		-v ${PWD}:/local openapitools/openapi-generator-cli generate \
    	-i /local/docs/internal-api.yml \
    	-g go \
    	-o /local/generated/go/internal

