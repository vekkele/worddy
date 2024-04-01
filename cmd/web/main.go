package main

import (
	"flag"
	"log"
	"net/http"
)

type application struct {
	addr string
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	app := application{
		addr: *addr,
	}

	err := http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
