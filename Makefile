
USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

build: identity
#openapitools

identity:
	pkger
	go build -o ${PWD}/bin/identity cmd/identity/*

rotate:
	go build -o ${PWD}/bin/rotate cmd/rotate/*
	./bin/rotate

# Generate the go code from the public and internal api's
openapitools:
	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/.openapi-generator/client-generator-config.yml /local/.openapi-generator/server-generator-config.yml

run: identity
	./bin/identity

migrate:
	pkger
	cd ./cmd/migrate && go build -o $(PWD)/bin/migrate
	./bin/migrate

install:
	go get github.com/markbates/pkger/cmd/pkger

docker:
	docker build -f Dockerfile -t moov/identity .
	docker run -it --rm moov/identity

clean:
	rm ./data/*