package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/bkhtbkid/go_final_project/pkg/api"
	"github.com/bkhtbkid/go_final_project/pkg/db"
	"github.com/bkhtbkid/go_final_project/pkg/server"
	"github.com/joho/godotenv"
)

//go:embed web
var web embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using OS environment")
	}

	dbFile := os.Getenv("TODO_DBFILE")

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal("Error init db", err)
	}

	port := os.Getenv("TODO_PORT")

	http.Handle("/", server.FsHandler(web))
	api.Init()

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
