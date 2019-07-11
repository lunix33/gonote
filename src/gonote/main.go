package main

import (
	"gonote/route"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Register the web routes.
	route.RegisterRoute()

	// HTTP listen on port 8080
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
