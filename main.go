package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"kirgu.ru/employee/repository"
	"kirgu.ru/employee/server"
)

func main() {

	prep, err := repository.NewPostgresRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	server := server.NewServer(e, prep)
	server.Start()
}
