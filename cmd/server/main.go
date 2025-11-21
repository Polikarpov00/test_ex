package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/go-qa-api/internal/app"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := app.NewGormDB()
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	app.AutoMigrate(db)

	a := &app.App{DB: db}

	log.Printf("server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, a.Routes()))
}
