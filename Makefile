DC=docker compose -f ./docker-compose.yml -p request-sink --env-file .dockerenv

.PHONY: dcup dcdown prep-env exec dc tidy tail restart fmt scan-fs update-deps

prep-env:
	@if [ ! -f .dockerenv ]; then \
		cp .dockerenv.example .dockerenv; \
	fi

dcup: prep-env
	$(DC) up -d

dcdown: prep-env
	$(DC) down --remove-orphans

exec:
	$(DC) exec app bash

tidy: 
	$(DC) exec app go mod tidy

fmt: 
	$(DC) exec app go fmt ./...

dc:
	@echo $(DC)

tail:
	$(DC) logs -f app

scan-fs:
	$(DC) run --rm trivy-scan fs .

update-deps:
	$(DC) exec app go get -u ./...

restart: dcdown dcup