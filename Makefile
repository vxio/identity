
USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

GEN_CODE_LOCATION := "pkg/gen"

build: compile
#openapitools

pkger:
	pkger

compile: pkger
	cd ./cmd/identity && go build -o $(PWD)/bin/identity
	cd ./cmd/rotate && go build -o ${PWD}/bin/rotate

# Generate the go code from the public and internal api's
openapitools:
	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/.openapi-generator/client-generator-config.yml /local/.openapi-generator/server-generator-config.yml

run: compile
	rm ./bin/identity.db
	./bin/identity

rotate:	compile
	./bin/rotate

migrate:
	pkger
	cd ./cmd/migrate && go build -o $(PWD)/bin/migrate
	./bin/migrate