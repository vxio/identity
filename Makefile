
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

test: build
	docker-compose up -d
	go test -cover ./...

migrate:
	pkger
	cd ./cmd/migrate && go build -o $(PWD)/bin/migrate
	./bin/migrate

install:
	go get github.com/markbates/pkger/cmd/pkger

docker: install test
	pkger
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ${PWD}/bin/.docker/identity cmd/identity/*
	docker build -f Dockerfile -t moov/identity .
	
docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov/identity

clean:
	rm ./data/*

# Generate the go code from the public and internal api's
openapitools:
	docker run --rm \
		-u $(USERID):$(GROUPID) \
		-e OPENAPI_GENERATOR_VERSION='4.2.0' \
		-v ${PWD}:/local openapitools/openapi-generator-cli batch -- /local/.openapi-generator/client-generator-config.yml
