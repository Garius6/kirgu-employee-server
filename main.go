package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()
	repo := NewMapRepository()
	server := NewServer(e, repo)
	server.Start()
}
