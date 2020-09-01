PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

# General make commands for projects

.PHONY: build
build: identity

.PHONY: check
check: build services
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	./lint-project.sh
endif

dist: clean build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/identity.exe cmd/identity/*
else
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -o bin/identity-$(PLATFORM)-amd64 cmd/identity/*
endif

docker: clean install
	pkger
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ${PWD}/bin/.docker/identity cmd/identity/*

# Docker image
	docker build --pull -t moov/identity:$(VERSION) -f Dockerfile .
	docker tag moov/identity:$(VERSION) moov/identity:latest

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/
endif

.PHONY: cover-test
cover-test: services
	go test -coverprofile=cover.out ./...

.PHONY: cover-web
cover-web: services
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

release: docker AUTHORS
	go vet ./...
	go test -coverprofile=cover-$(VERSION).out ./...
	git tag -f $(VERSION)

docker-push:
	docker push moov/identity:$(VERSION)
	docker push moov/identity:latest

quay-push:
	docker push quay.io/moov/identity:$(VERSION)
	docker push quay.io/moov/identity:latest

# Custom to go-services

docker-run: docker
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov/identity:$(VERSION)

install:
	go install github.com/markbates/pkger/cmd/pkger
	git checkout LICENSE

identity:
	pkger
	go build -o ${PWD}/bin/identity cmd/identity/*

run: identity
	./bin/identity

test: services build
	go test -cover ./...

services:
	-docker-compose up -d --force-recreate

# Generate the go code from the public and internal api's
openapitools:
	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 batch -- /local/.openapi-generator/client-generator-config.yml
	gofmt -w ./pkg/client/


rotate:
	go build -o ${PWD}/bin/rotate cmd/rotate/*