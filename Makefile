PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

build: identity rotate

identity:
	pkger
	go build -o ${PWD}/bin/identity cmd/identity/*

rotate:
	go build -o ${PWD}/bin/rotate cmd/rotate/*

run: identity
	./bin/identity

test: services build
	go test -cover ./...

services:
	-docker-compose up -d

migrate:
	pkger
	cd ./cmd/migrate && go build -o $(PWD)/bin/migrate
	./bin/migrate

install:
	go get github.com/markbates/pkger/cmd/pkger
	git checkout LICENSE

.PHONY: check
check: build services
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	./lint-project.sh
endif

docker: install
	pkger
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ${PWD}/bin/.docker/identity cmd/identity/*
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
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/.openapi-generator/client-generator-config.yml

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@