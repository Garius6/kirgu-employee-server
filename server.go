package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"kirgu.ru/employee/model"
)

type Repository interface {
	SignIn(username string, password string) (*model.User, error)
	SignUp(username string, password string, passwordConfirmation string) error
}

type Server struct {
	e    *echo.Echo
	repo Repository
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func NewServer(e *echo.Echo, repo Repository) *Server {
	return &Server{e, repo}
}

func (s *Server) Start() {
	s.e.POST("/users/sign_in", s.SignIn)
	s.e.POST("/users/sign_up", s.SignUp)
	log.Fatal(s.e.Start(":8090"))
}

func (s *Server) SignIn(c echo.Context) error {

	var signInRequest SignInRequest
	err := c.Bind(&signInRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	user, err := s.repo.SignIn(signInRequest.Username, signInRequest.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) SignUp(c echo.Context) error {

	var signUpRequest SignUpRequest
	err := c.Bind(&signUpRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = s.repo.SignUp(signUpRequest.Username, signUpRequest.Password, signUpRequest.PasswordConfirmation)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "Created")
}
