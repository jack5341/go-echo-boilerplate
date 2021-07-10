package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Username string
	Password string
}

func Register(e *echo.Group) {
	e.POST("/register", func(c echo.Context) error {

		user := User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Password == " " {
			return c.String(http.StatusBadRequest, "password is invalid !")
		}

		if user.Username == "" {
			return c.String(http.StatusBadRequest, "username is invalid !")
		}
		return c.JSON(http.StatusOK, user)
	})
}

func Login() {

}
