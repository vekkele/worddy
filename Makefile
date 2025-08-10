-include .env
export

.PHONY: css
css:
	./tailwindcss -i ./ui/tailwind/main.css -o ./ui/static/css/app.css --watch

tailwindcss:
	curl -fL -o ./tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.11/tailwindcss-$(TAILWINDCSS_OS_ARCH)
	chmod a+x ./tailwindcss

.PHONY: deps
deps:
	go mod download

.PHONY: init
init: deps install 

.PHONY: build-css
build-css: tailwindcss
	./tailwindcss -i ./ui/tailwind/main.css -o ./ui/static/css/app.css --minify

.PHONY: build-prod
build-prod: build-css
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/main ./cmd/web

.PHONY: install
install: tailwindcss
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: db_init
db_init: 
	docker compose --env-file .env --file ./db/dev/docker-compose.yml up -d

.PHONY: db_rm
db_rm:
	docker compose --file ./db/dev/docker-compose.yml down -v

.PHONY: db_migration
db_migration:
	goose -dir internal/store/postgres/migrations -s create "$(name)" sql

.PHONY: db_status
db_status:
	goose -dir internal/store/postgres/migrations postgres $(WORDDY_DB_DSN) status

.PHONY: db_up
db_up:
	goose -dir internal/store/postgres/migrations postgres $(WORDDY_DB_DSN) up

.PHONY: db_down
db_down:
	goose -dir internal/store/postgres/migrations postgres $(WORDDY_DB_DSN) down

.PHONY: db_connect
db_connect:
	docker exec -it worddy_db psql -d $(WORDDY_DB_DSN)

.PHONY: cover
cover:
	go test -coverprofile=./tmp/coverage.out ./...
	go tool cover -html=./tmp/coverage.out
