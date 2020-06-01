PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

build: identity

identity:
	./pkger
	go build -o ${PWD}/bin/identity cmd/identity/*

rotate:
	go build -o ${PWD}/bin/rotate cmd/rotate/*
	./bin/rotate

run: identity
	./bin/identity

migrate:
	./pkger
	cd ./cmd/migrate && go build -o $(PWD)/bin/migrate
	./bin/migrate

install:
	wget -O pkger.tar.gz https://github.com/markbates/pkger/releases/download/v0.16.0/pkger_0.16.0_$(shell uname -s)_x86_64.tar.gz
	tar xf pkger.tar.gz

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	time ./lint-project.sh
endif

docker: install
	./pkger
	docker build --pull -t moov/identity:$(VERSION) -f Dockerfile .
	docker tag moov/identity:$(VERSION) moov/identity:latest

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov/identity:$(VERSION)

clean:
	rm ./data/*

# Generate the go code from the public and internal api's
openapitools:
	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/.openapi-generator/client-generator-config.yml /local/.openapi-generator/server-generator-config.yml
