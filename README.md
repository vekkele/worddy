# Worddy

App to learn new vocabulary using SRS written in Go + HTMX + TailwindCSS

## Development
1. Create `.env` file, use `.env.example` as reference
2. Run `make init` to intall dependencies
3. Run `make db_init` to create dev db in docker container (make sure docker is installed)
4. Run `make db_up` to run all db migrations
5. Run `make css` to run tailwind in watch mode
6. Run `air` to start dev server