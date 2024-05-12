package server

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Username,
	})

	t, err := token.SignedString(secret)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
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
