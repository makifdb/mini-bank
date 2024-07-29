package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makifdb/mini-bank/minimalist/internal/app"
	"github.com/makifdb/mini-bank/minimalist/internal/redis"
	"github.com/makifdb/mini-bank/minimalist/internal/repository"
	"github.com/makifdb/mini-bank/minimalist/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := repository.NewDatabase(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mailService := utils.NewMailService()

	rds, err := redis.NewClient(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	router := app.NewApp(db, rds, mailService)

	log.Println("Server is running on port", os.Getenv("PORT"))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatal(err)
	}
}
