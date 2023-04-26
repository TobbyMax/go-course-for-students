package main

import (
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	"homework9/internal/ports/httpgin"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New()))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
