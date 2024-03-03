package main

import "flag"

type application struct {
	addr string
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	app := application{
		addr: *addr,
	}

	app.start()
}
