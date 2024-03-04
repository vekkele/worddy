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