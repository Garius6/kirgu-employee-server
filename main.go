package main

import (
	"github.com/labstack/echo/v4"
	"kirgu.ru/employee/server"
)

func main() {
	e := echo.New()
	repo := NewMapRepository()
	server := server.NewServer(e, repo)
	server.Start()
}
