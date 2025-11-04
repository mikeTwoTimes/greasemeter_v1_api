package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/app"
)

func main() {
	a, err := app.NewApp(
		os.Getenv("PORT"),
		os.Getenv("DB_CONN"),
		os.Getenv("JWT_SECRET"),
		os.Getenv("SENDGRID_KEY"),
	)

	if err != nil {
		log.Fatal(err)
	} else if err = a.Serve(); err != nil {
		log.Fatal(err)
	}
}
