SHELL=/bin/bash

UID := $(shell id -u)

up:
	env UID=${UID} docker-compose up -d --build --remove-orphans --force-recreate

stop:
	env UID=${UID} docker-compose stop

remove:
	env UID=${UID} docker-compose rm

bash:
	env UID=${UID} docker-compose exec -u app php-ps-fpm bash

network:
	docker network create --attachable local
