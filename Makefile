include .env
export

css:
	./tailwindcss -i ./ui/tailwind/main.css -o ./ui/static/css/app.css --watch

tailwind_install:
	rm -f tailwindcss
	curl -L https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64 -o tailwindcss
	chmod +x tailwindcss

deps:
	go mod download

init: deps install 

install: tailwind_install
	go install github.com/cosmtrek/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

db_init: 
	docker compose --env-file .env --file ./db/dev/docker-compose.yml up -d

db_rm:
	docker compose --file ./db/dev/docker-compose.yml down -v

db_migration:
	goose -dir db/migrations -s create "$(name)" sql

db_status:
	goose -dir db/migrations postgres $(WORDDY_DB_DSN) status

db_up:
	goose -dir db/migrations postgres $(WORDDY_DB_DSN) up

db_down:
	goose -dir db/migrations postgres $(WORDDY_DB_DSN) down

cover:
	go test -coverprofile=./tmp/coverage.out ./...
	go tool cover -html=./tmp/coverage.out
