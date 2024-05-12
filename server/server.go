package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"kirgu.ru/employee/model"
)

var (
	secret = []byte("secret")
)

type Repository interface {
	SignIn(username string, password string) (*model.User, error)
	SignUp(username string, password string, passwordConfirmation string) error
}

type Server struct {
	e    *echo.Echo
	repo Repository
}

func NewServer(e *echo.Echo, repo Repository) *Server {
	return &Server{e, repo}
}

func (s *Server) Start() {

	s.e.POST("/users/sign_in", s.SignIn)
	s.e.POST("/users/sign_up", s.SignUp)

	log.Fatal(s.e.Start(":8090"))
}
