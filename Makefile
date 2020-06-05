SHELL=/bin/bash

UID := $(shell id -u)

get-env:
	cp -a ./.env.example.common ./.env
up:
	docker-compose up -d --build --remove-orphans --force-recreate

stop:
	env UID=${UID} docker-compose stop

remove:
	env UID=${UID} docker-compose rm

bash:
	docker-compose exec redis-dev bash

network:
	docker network create --attachable local

# build targets
build: build-installer build-grabber build-parser

build-installer:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/installer ./cmd/cfgInstaller.go ./cmd/CommonVars.go

build-grabber:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/grabber ./cmd/grabber.go ./cmd/CommonVars.go

build-parser:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/parser ./cmd/parser.go ./cmd/CommonVars.go

build-env-file:
	./bin/installer -p=./bin

test:
	go test -v ./cmd/*test.go

test-adv:
	go test -v -race ./cmd/*test.go
